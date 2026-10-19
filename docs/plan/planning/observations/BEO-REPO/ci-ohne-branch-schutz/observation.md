# Die CI ist ohne Branch-Schutz ein Hinweis, kein Sperrriegel

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`.github/workflows/ci.yml` fährt `make gates` und die Prüfung der Commit-Range auf GitHub; ein roter Lauf färbt den Status des Jobs `ci-gates`, hält aber keinen Merge auf, solange auf GitHub kein „required status check“ für diesen Namen gesetzt ist. Die Einstellung liegt auf GitHub, ist eine Handlung des Betreibers und aus dem Klon nicht prüfbar; das Repo kann sie nennen (`harness/README.md` §Sensors nennt den Namen `ci-gates` und die Handlung), aber nicht behaupten, dass sie gesetzt ist (`AGENTS.md` §3.11). Bis sie gesetzt ist, schützt die CI nur den, der den Status liest.

Abgrenzung zu `BEO-REPO/range-pruefung-laeuft-nicht-im-worktree`: Dort läuft die Prüfung lokal nicht; hier läuft sie auf GitHub und bindet nichts.
