**Stand:** offen

Zweites Auftreten (`slice-release-v0-2-0`, Urteil des Planners, 2026-10-24; erstes `slice-erster-release-homebrew-nachweis`, 2026-10-21): beide Slices lieferten den Nachweis für Linux/`amd64` und führen macOS und `arm64` als Grenze. Ohne Adresse: `slice-erster-release-homebrew-nachweis` liegt in `done/` und nimmt nichts an (`AGENTS.md` §3.13). Ein Sensor ist im Klon nicht möglich (kein macOS und kein Apple-Silicon-Host unter den Voraussetzungen `bash`, `git`, `docker`).

**Vorschlag des Planners, nicht entschieden:** (a) Die Grenze bleibt dauerhaft in `LH-FA-19.a`, weil der Beleg ein Lauf des Betreibers auf einem Mac wäre; dann `gestrichen` mit dieser Begründung, wenn der Nutzer es so entscheidet. (b) Ein offener Slice, dessen Handlung der Nutzer auf einem Mac fährt (Tap einbinden, Paket installieren, `version`, Paket-Test), mit der Ausgabe als Beleg-Datei nach `SPEC-067` Punkt 20 (6); der Planner legt ihn nur an, wenn der Nutzer es will. Ausgang nicht zugewiesen.
