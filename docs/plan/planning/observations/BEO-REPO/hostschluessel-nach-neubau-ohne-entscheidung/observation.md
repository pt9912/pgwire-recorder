# Nach Neubau des Images warnt SSH am Host, und die Entscheidung über eine Einstellung liegt offen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Dev-Container tauscht Stände über SSH aus. Nach einem Neubau des Images ohne Cache hat der Container andere Hostschlüssel; `git fetch` und `git push` am Host enden mit `REMOTE HOST IDENTIFICATION HAS CHANGED` und Exit 128 bzw. 255. Ob der Betrieb dafür eine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) nutzt, ist Entscheidung des Nutzers und offen; `docs/developer/devcontainer.md` nennt, was gelaufen ist, und empfiehlt keine Einstellung (§13).
