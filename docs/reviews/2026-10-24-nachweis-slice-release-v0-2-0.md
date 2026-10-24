# Nachweis: Installation aus dem Tap für `v0.2.0` unter Linux auf `x86_64`

**Datum:** 2026-10-24, Beginn 08:06:34Z, Ende 08:06:49Z (Zeilen `date -u` der Ausgabe). **Stand:** Commit `90199364ecfbd3545eb27094f911face6e1a5008` (Zweig `impl-rel-v020c`, Arbeitsbaum `/Development/pgwire-recorder-rel5`, sauber); `tools/harness/homebrew-nachweis.sh` ist an diesem Stand gleich dem von `main` bei `511428c9`. **Plattform:** Linux, `x86_64` (`uname -s`, `uname -m` im Lauf, im Container). **Tag und Digest:** Release `v0.2.0` auf dem Commit `32d0d84a924a1e67a44e9cee9441a69b22eb4fdd` (Tag-Objekt `fa5de1591785171ba4305dcd063343569a5d96ff`); Image des Laufs `homebrew/brew:4.6.20@sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd`; Index-Digest des Produkt-Image `0.2.0` `sha256:1d50181235d6a0317aaf376c33c1f17b421fb25036dea3324946a652dd2ceee4` (aus Lauf `40074219261`, kein Teil dieses Laufs).

Aufruf: `make homebrew-nachweis` ([`SPEC-067`](../../harness/vertraege/SPEC-067-release.md) Punkt 20, Abnahmeszenario 11, Anforderung [`LH-FA-19`](../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)), Exit-Code von `make` 0 (letzte Zeilen des Codeblocks). Befehle (a) bis (e): (a) `brew tap`, (b) `brew install`, (c) `pgwire-recorder version`, (d) `brew test`, (e) Replay-Lauf mit `SIGTERM`; der Exit-Code des Replay-Laufs ist 0. Dies ist ein zweiter Lauf: die volle Ausgabe des ersten Laufs (Beginn 07:42:09Z, Eintrag `freigaben/0.2.0.md`, Punkt 11) wurde nicht gesichert und zählt nicht als Beleg.

**Gekürzt** ist genau eine Stelle: in der Zeile `==> Downloading from https://release-assets.githubusercontent.com/…` die Abfrage der signierten Adresse (Signatur und Ablaufzeit), Marke `[gekürzt: signierte Adresse]`; Host, Pfad und der Rest des Laufs stehen vollständig. Gekürzt wurde beim Ablegen der Datei, vor dem Commit.

**Was die Ausgabe für die Sätze des Handbuchs zeigt** (gelesen am Text darunter, kein weiterer Lauf): `brew install` holt eine Datei (`Fetching downloads for: pgwire-recorder`, ein `Downloading` von `…/v0.2.0/pgwire-recorder-linux-amd64`), nennt keine Abhängigkeit und keinen Übersetzungsschritt; die Zeile `built in 3 seconds` ist der feste Text von Homebrew für das Ende des Installationsschritts. Installiert ist `…/Cellar/pgwire-recorder/0.2.0`; die Tags `v*` des Repositorys sind zu diesem Zeitpunkt `v0.1.0`, `v0.1.1` und `v0.2.0` (`git ls-remote --tags origin 'v*'`), `0.2.0` ist die höchste stabile Version. Die Ausgabe zeigt nicht, wie die Formel die Version wählt; die Aussage „trägt die höchste stabile Version“ gilt damit für den Stand dieses Laufs.

````text
homebrew-nachweis: Image homebrew/brew:4.6.20@sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd
+ date -u +%Y-%m-%dT%H:%M:%SZ
2026-10-24T08:06:34Z
+ uname -s
Linux
+ uname -m
x86_64
+ brew tap pt9912/pgwire-recorder
==> Tapping pt9912/pgwire-recorder
Cloning into '/home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/pt9912/homebrew-pgwire-recorder'...
Tapped 1 formula (13 files, 69.1KB).
+ brew install pt9912/pgwire-recorder/pgwire-recorder
==> Fetching downloads for: pgwire-recorder
==> Fetching pt9912/pgwire-recorder/pgwire-recorder
==> Downloading https://github.com/pt9912/pgwire-recorder/releases/download/v0.2.0/pgwire-recorder-linux-amd64
==> Downloading from https://release-assets.githubusercontent.com/github-production-release-asset/1403356322/22aab9ed-ccf6-493c-b6ab-08ee01af2eba?[gekürzt: signierte Adresse]
==> Installing pgwire-recorder from pt9912/pgwire-recorder
🍺  /home/linuxbrew/.linuxbrew/Cellar/pgwire-recorder/0.2.0: 4 files, 11.4MB, built in 3 seconds
==> Running `brew cleanup pgwire-recorder`...
+ pgwire-recorder version
pgwire-recorder 0.2.0
+ brew test pgwire-recorder
Warning: test is a developer command, so Homebrew's
developer mode has been automatically turned on.
To turn developer mode off, run:
  brew developer off

✔︎ JSON API cask_tap_migrations.jws.json
==> Testing pt9912/pgwire-recorder/pgwire-recorder
==> /home/linuxbrew/.linuxbrew/Cellar/pgwire-recorder/0.2.0/bin/pgwire-recorder version
+ pgwire-recorder replay --listen 127.0.0.1:55432 --input /tmp/tmp.VmKNshhdAu/aufzeichnung.yaml (Fehlerausgabe in eine Datei)
+ warten, bis der Replay-Lauf seine erste Zeile geschrieben hat (Groesse der Fehlerausgabe, Abstand 1 s, Grenze 60 s)
+ kill -TERM <Replay-Lauf>
--- Fehlerausgabe des Replay-Laufs ---
time=2026-10-24T08:06:48.108Z level=INFO msg="replay gestartet" listen=127.0.0.1:55432 input=/tmp/tmp.VmKNshhdAu/aufzeichnung.yaml
time=2026-10-24T08:06:49.106Z level=INFO msg="Herunterfahren begonnen" sessions=0
time=2026-10-24T08:06:49.106Z level=WARN msg="1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 1" code=PGR-W2001
time=2026-10-24T08:06:49.106Z level=INFO msg="replay beendet"
--- Exit-Code des Replay-Laufs: 0
homebrew-nachweis: Skript endet mit Exit-Code 0
homebrew-nachweis: der Lauf endet mit Exit-Code 0
make-exit=0
2026-10-24T08:06:49Z
````
