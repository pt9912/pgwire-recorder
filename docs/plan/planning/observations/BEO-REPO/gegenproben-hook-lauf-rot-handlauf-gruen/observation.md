# Der Hook-Lauf der Gegenproben im Gegenproben-Container ist rot, der Handlauf auf demselben Hash grün

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der `post-receive`-Hook im Bare-Repo des Gegenproben-Containers startet `make gegenproben` im Hintergrund auf dem eingelieferten Hash. Der Lauf endet rot, wo derselbe Hash von Hand per `ssh` im selben Container grün endet. Beobachtet: Der Hintergrundlauf ignorierte `SIGINT`, `a-check` meldete „permission denied“ auf `/src`, acht Gegenproben waren rot. Die Ursache ist nicht geklärt; Kandidaten sind die Umgebung des Hintergrundlaufs (Signal-Disposition, Benutzer, Dateirechte), und es ist offen, ob es dieselbe Klasse ist wie die zwei Nachbarn.

Abgrenzung zu `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale`: Dort sind Deskriptoren und ignorierte Signale als Ursache gemessen und im Hook behoben; hier steht die Ursache offen und der Lauf ist nach der Behebung rot. Abgrenzung zu `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`: Dort hängt ein Fall an einer Einstellung des Hosts; hier unterscheidet sich der Aufrufer, nicht der Host. Die Zuordnung zu einer der beiden folgt der Ursache.
