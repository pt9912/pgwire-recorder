# Welche Laufzeit der Dev-Container für verschachtelte Container nutzt, ist offen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`u-boot` 0.10.0 kennt für `nestedRuntime` die Werte `docker` und `podman`. Gemessen ist `docker` (Daemon im Container, `make gates` grün). Für `podman` nennt der Changelog von `u-boot` `podman build` auf Colima (arm64) mit `/dev/net/tun`, ohne `buildx`, `docker compose` und `# syntax=`; das `Dockerfile` des Repos nutzt `RUN --network=none`, `--platform=$BUILDPLATFORM` und mehrstufige Builds, ob `make` mit `podman` läuft, ist ungeprüft. Die Wahl braucht den Nutzer; entscheidet er für `podman`, legt der Planner einen Folge-Slice an. `docs/developer/devcontainer.md` §13 führt die Frage ohne Empfehlung.
