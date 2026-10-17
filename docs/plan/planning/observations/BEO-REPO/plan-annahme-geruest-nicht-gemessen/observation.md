# Ein Plan nennt Aufrufe und Wirkungen eines Fremdwerkzeugs, die der Slice erst misst, und ein Nehmer übernimmt sie

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Plan eines Gebers beschreibt den Aufruf oder die Wirkung eines Werkzeugs außerhalb des Repos (`u-boot up` und `u-boot down`, Port-Abbildung in `compose.yaml`, `remoteEnv` „beim Start“) nach dem Auftrag, nicht nach einer Messung. Ein Nehmer übernimmt die Sendung im selben Wortlaut (Pfade, Neustart, Übergabe des Zugangs); die Messung des Gebers widerlegt sie, und der Nehmer läuft nach seinem Text ins Leere. Abgrenzung zu `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`: Die Adresse nimmt an, aber mit einem Inhalt, der nicht trägt; und zu `BEO-REPO/plan-folgt-korrektur-nicht`: Hier folgt der Plan nicht einem Diff, sondern der Nehmer einer Messung des Gebers.
