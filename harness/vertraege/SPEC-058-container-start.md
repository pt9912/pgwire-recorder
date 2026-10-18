# SPEC-058 — Start des Dev-Containers (`container-start`, `container-start-gegenprobe`)

`make container-start` legt am Host das Gerüst des Dev-Containers an, wenn es fehlt, und startet den
Container. `make container-start-gegenprobe` prüft das Werkzeug. Das Werkzeug ist kein Gate: Es steht
im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`, und kein Hook und kein Stand
glaubt sein Ergebnis. Diese Kennung legt fest, was es zusagt, woran es scheitert und was die Gegenprobe
belegt. Sie sagt nichts über das, was im Container liegt (Bare-Repo, Klon, Zweig, Zugang von Claude Code),
und nichts über `docs/developer/devcontainer.md`, außer wo ein Punkt es nennt.

Der Vertrag hat zwei Teile, damit ein Schnitt an einer Naht liegt: **Teil I** (Punkte 2 bis 8) prüft Eingaben,
Voraussetzungen und Zustand und legt das Gerüst an, ohne dass Docker etwas ändert; **Teil II** (Punkte 9 bis 12)
startet den Container. Punkte 13 bis 17 gelten für beide.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element
einer Aufzählung gibt es einen Fall und eine Mutation, die nur ihn rot färbt; die Gegenprobe verlangt
genau diese Menge roter Fälle (mehr und weniger machen sie rot).

**Plattform.** Zugesagt ist ein Linux-Host mit bash ≥ 4 und GNU-Werkzeugen, gemessen mit Docker 29.9.1,
`u-boot` 0.10.0 und `devcontainer` 0.80.2. macOS- und Windows-Hosts, Docker Desktop und andere Versionen
von `devcontainer` sind nicht zugesagt (akzeptiertes Negativ; Grund: kein Host dieser Art steht zur
Messung bereit, und das Ziel lehnt sie nicht ab, damit ein Nutzer es dort auf eigene Gefahr versuchen kann).

**Randform → Punkt.**

| Randform des Plans | Punkt |
|---|---|
| A1 Werkzeug fehlt oder falsche Version | 3 |
| A2 `PGR_DEV` und die übrigen Eingaben | 2, 5 (Lage) |
| A3 Verzeichnis besteht schon | 5 |
| A4 Port | 2 (Form), 10 (belegt) |
| A5 zweiter Aufruf, Name, Kollision | 9 |
| A6 Schlüsselpaar | 6 |
| A7 Warnungen und Fehler von `u-boot` | 7 |
| A8 Bereitschaft, Zeitgrenze | 11 |
| A9 Hostschlüssel | 11 |
| A10 Abbruch | 12 |
| A11 Reihenfolge der Gründe | 13 |
| A12 Stopp und Neustart | 9 |
| Aufrufarten von `make` | 14 |
| Plattform | Kopf dieser Kennung, 17 |

1. **Gegenstand und Ort.** Das Ziel `container-start` ist im Repo (Fragment `harness/mk/container-start.mk`,
   Skript `tools/harness/container-start.sh`), `.PHONY`, kein Gate. Es schreibt nur unter `$PGR_DEV`
   (Punkt 2) und startet Container; es berührt das Host-Repo nicht (kein `git`-Aufruf mit Schreibwirkung,
   kein Remote, kein `core.sshCommand`, keine Datei im Arbeitsbaum), kein `~/.ssh` und kein `known_hosts`.
   Es entfernt, stoppt und baut nie etwas ab: kein `docker stop`, `docker rm`, `docker rmi`,
   `docker volume rm`, `prune`. Die Gegenprobe zeigt es an einer Attrappe von `docker`, die jeden Aufruf
   protokolliert (Fall `p1-keine-zerstoerung`; Mutation in Teil I: ein Schritt von Teil I, etwa `werkzeuge`, ruft `docker rm -f`; mit Teil II zusätzlich: das Skript ruft bei Kollision (Punkt 9) `docker rm -f`).
   Stoppen ist `docker stop <Name>` von Hand; ein zweites Ziel gibt es nicht (Punkt 9).

2. **Eingaben.** Sie kommen aus der Umgebung. Ein Wert, der gesetzt und leer ist, ist ungültig (kein
   stilles Einsetzen der Vorgabe); nicht gesetzt heißt Vorgabe.

   | Variable | Vorgabe | gültig |
   |---|---|---|
   | `PGR_DEV` | keine, Pflicht | absoluter Pfad aus `A-Za-z0-9._+-` und `/`; kein Segment `.`, `..` oder leer; kein `/` am Ende; Basisname (der Projektname) aus `a-z0-9-`, erstes Zeichen `a-z0-9` |
   | `PGR_NAME` | `pgr-dev` | `[a-z0-9][a-z0-9_.-]*` |
   | `PGR_SSH_PORT` | `2222` | ganze Zahl 1 bis 65535, Dezimalziffern, keine führende Null, kein Vorzeichen |
   | `PGR_WARTEN` | `60` | ganze Zahl 1 bis 600, Form wie der Port (Sekunden, Punkt 11) |
   | `PGR_CPUS`, `PGR_MEM` | `8`, `8g` | nicht geprüft: unverändert an `mk-override.sh` gereicht (Punkt 8); ein ungültiger Wert zeigt sich als Fehler von Docker in Schritt `start` (Punkt 17, Grenze) |

   **Grund der Vorgaben** (Entscheidung des Architect nach Empfehlung des Planners; der Nutzer legte am 2026-10-18 fest, dass die Limits niedriger
   werden und das Ziel sie anpasst). Gemessen: `make gates` im Container mit Limit `12g` brauchte höchstens 5,42 GiB (zweiter Lauf 4,24 GiB,
   `docker stats`, alle 10 s), `OOMKilled=false`; im Leerlauf 0,2 bis 0,5 GiB; Host 31 GiB Speicher und 20 CPUs. Ein Limit reserviert nichts: es
   schützt den Host (etwa vor einer Rekursion) und macht die Überbuchung bei drei bis sechs Containern planbar. `8g` liegt rund 48 % über der
   gemessenen Spitze; `12g` lag bei mehr als dem Doppelten und erlaubte bei drei Containern 36 GiB auf einem Host mit 31 GiB. `8` CPUs bleiben:
   nur diese Zahl ist gemessen (Lauf von 1701 s), ein kleinerer Wert verlängerte den Lauf um ein Maß, das niemand kennt, und das CPU-Limit ist
   eine Quote, die ruhende Container nicht beansprucht. **Nicht gemessen:** ein Lauf von `make gates` bei `8g` mit Zeitmessung; mehrere
   Container zugleich unter Last. Wer mehr Container zugleich unter Last fährt, senkt `PGR_CPUS` oder `PGR_MEM` für sie.
   **Ändern ohne Neustart:** Ein laufender Container nimmt neue Werte mit `docker update --memory <x> --memory-swap <x> <Name>` an (beide Werte
   gleich, sonst bleibt der Swap des alten Wertes); das ist Hinweis im Dokument `docs/developer/devcontainer.md`, nicht Teil des Ziels
   (Punkt 9).

   Schritt `eingaben` endet bei einem Verstoß mit Exit-Code 2, ändert nichts und ruft kein Werkzeug auf. Er
   meldet **alle** Verstöße, je einer in einer Zeile, die die Variable nennt, in der Reihenfolge der Tabelle.
   Je Form ein Fall und eine Mutation, die nur ihre Prüfung entfallen lässt: `PGR_DEV` nicht gesetzt, leer,
   relativ, mit Leerzeichen, mit Tabulator, mit Zeilenumbruch, mit Nicht-ASCII-Zeichen, mit Segment `..`,
   mit Segment `.`, mit leerem Segment (`//`), mit `/` am Ende, mit Basisname aus Großbuchstaben, mit
   Basisname mit Unterstrich; `PGR_NAME` leer und mit Großbuchstaben; `PGR_SSH_PORT` leer, `abc`, `02222`,
   `+22`, `0`, `65536`; `PGR_WARTEN` leer, `abc`, `0`, `601`. Grenzen angenommen: `PGR_SSH_PORT=1` und
   `65535`, `PGR_WARTEN=1` und `600`. Nicht gesetzt ergibt die Vorgabe (Fall `p2-vorgaben`, sichtbar in
   `override.json`, Punkt 8; Mutation: andere Vorgabe). Fälle `p2-*`.

3. **Voraussetzungen am Host.** Schritt `werkzeuge`: `docker`, `u-boot`, `devcontainer`, `node`, `ssh`,
   `ssh-keygen` müssen im Pfad liegen. Dazu: (a) `u-boot --version` nennt die erste Zeichenfolge
   `<x>.<y>.<z>` aus Dezimalziffern; sie ist, Teil für Teil als Zahl verglichen, nicht kleiner als `0.10.0`
   (`0.9.0` ist kleiner, `0.100.0` ist nicht kleiner); ohne eine solche Zeichenfolge ist die Version
   unlesbar und abgelehnt (Meldung `u-boot: Version unlesbar`); `u-boot --version` mit Exit-Code ungleich 0 ist ebenso unlesbar, gleich was es ausgab; (b) `docker info` endet mit Exit-Code 0, und die Hauptversion des Servers
   (`docker version --format '{{.Server.Version}}'`) ist nicht kleiner als 24. Liefert `docker version --format '{{.Server.Version}}'` Exit-Code ungleich 0 oder keine mit Dezimalziffern beginnende Hauptversion, ist die Version unlesbar und abgelehnt (Meldung `docker: Version des Servers unlesbar`, verschieden von „kleiner als 24“; Exit 3). Version von `devcontainer`,
   `node`, `ssh`: nicht geprüft (Punkt 17). Ein Verstoß endet mit Exit-Code 3 und ändert nichts. Der Schritt meldet **alle**
   Verstöße, je einer in einer Zeile mit dem Namen des Werkzeugs, in der Reihenfolge `docker`, `u-boot`,
   `devcontainer`, `node`, `ssh`, `ssh-keygen`; bei einem fehlenden Werkzeug entfällt dessen Versionsprüfung.
   Fälle `p3-*`: je Werkzeug fehlt (sechs), `u-boot` 0.9.0 abgelehnt, 0.10.0, 0.100.0 und 1.0.0 angenommen,
   unlesbare Version abgelehnt (auch `u-boot --version` mit Exit 1, auch `docker version` ohne Ziffern oder mit Exit 1), Daemon nicht erreichbar, Server 23 abgelehnt und 24 angenommen, alle fehlen
   (Reihenfolge). Mutationen: das Werkzeug wird nicht geprüft, lexikalischer statt numerischer Vergleich, `>`
   statt `>=`, Docker-Version nicht geprüft, Reihenfolge vertauscht. **E2:** Diese Voraussetzungen gelten nur
   für `container-start` und seine Gegenprobe; kein anderes Ziel und kein Gate braucht sie. Stand: `ssh` ist
   Entscheidung des Nutzers (2026-10-17); `u-boot` ab 0.10.0, `devcontainer`, `node` und `ssh-keygen` sind nach
   Empfehlung entschieden, nicht vom Nutzer; streicht er eine, entfällt der Automatismus für das Gerüst.

4. **Ablauf.** Die Schritte laufen in dieser Reihenfolge, der erste Verstoß beendet den Lauf: `eingaben` (2),
   `werkzeuge` (3), `zustand` (5, 6 nur lesen, 9, 10), `geruest` (5, 7), `schluessel` (6), `override` (8),
   `start` (9, 10), `warten` (11). **Bis einschließlich `zustand` ändert das Ziel nichts:** weder Verzeichnis
   noch Datei noch Container, und es ruft von `u-boot` nur `--version` und von `docker` nur `info`,
   `version`, `ps`, `inspect` und `port` auf. Fall `p4-nichts-geaendert` je Verstoß der Schritte 2, 3, 5, 6, 9
   und 10: `$PGR_DEV` (Existenz, Inhalt, Rechte) ist unverändert, das Protokoll der Attrappen nennt keinen
   anderen Aufruf; Mutation: ein Schritt schreibt vor seiner Prüfung. Jede Fehlermeldung nennt den Schritt
   (`container-start: <Schritt>: <Grund>`, auf stderr); auf stdout steht von der eigenen Ausgabe des Ziels nur die Schlussmeldung des Erfolgs; die Ausgabe eines aufgerufenen Werkzeugs (`u-boot`, in Teil II `devcontainer` und die letzte Ausgabe von `ssh` in der Meldung) erscheint unverändert auf dem Strom, auf dem das Werkzeug sie schrieb, und ist davon ausgenommen (Punkt 7)
   (Punkt 11).

5. **Das Verzeichnis `$PGR_DEV`.** Zustand in Schritt `zustand`, Anlegen in `geruest`. Gleichwertig ist
   die Auflösung der Symlinks des vorhandenen Teils des Pfades; verglichen wird mit der Wurzel des Repos,
   in dem das Skript liegt (nicht mit dem Arbeitsverzeichnis).
   - Der Pfad ist die Wurzel oder liegt darunter (auch über einen Symlink): abgelehnt, Exit-Code 4. Ein Pfad, der nur
     textlich mit dem Namen der Wurzel beginnt (`<Wurzel>-zwei`), liegt nicht darunter und wird angenommen.
   - Der Elternpfad fehlt oder ist kein Verzeichnis: abgelehnt, Exit 4 (das Ziel legt nie mehr als ein
     Verzeichnis an, kein `mkdir -p`).
   - Der Pfad ist vorhanden und kein Verzeichnis: abgelehnt.
   - Fehlt das Verzeichnis oder ist es leer: es wird angelegt beziehungsweise genutzt; `u-boot init`
     läuft in ihm (Punkt 7).
   - Enthält es eine Datei `u-boot.yaml` mit der Zeile `    name: <Projektname>` und das Verzeichnis
     `.devcontainer/` mit `devcontainer.json`: **Gerüst vorhanden**. Es folgt kein Aufruf von `u-boot`
     außer `--version`; Dateien des Gerüsts (`u-boot.yaml`, `.devcontainer/`, `mk-override.sh`, `ssh/`) bleiben byteweise gleich, `.devcontainer/` ohne `override.json` (diese Datei schreibt Schritt `override` bei jedem Lauf neu, Punkt 8)
     (Fall: Prüfsummen vor und nach dem Lauf; Mutation: das Ziel ruft `u-boot generate devcontainer`).
   - Enthält es `u-boot.yaml` mit dem Namen, aber kein `.devcontainer/devcontainer.json`: ein früherer Lauf
     brach in `u-boot` ab; der Lauf setzt mit den `config set`-Zeilen und `generate` fort, ohne `init`.
   - Jeder andere Inhalt (Dateien ohne `u-boot.yaml`, `u-boot.yaml` mit anderem Namen) ist fremd: abgelehnt, Exit 4,
     die Meldung nennt das Verzeichnis; nichts wird überschrieben.
   Fälle `p5-*`, je Zeile einer; Mutationen: der Vergleich entfällt oder weicht ab (kein Auflösen der Symlinks,
   Präfixvergleich ohne `/`, `mkdir -p`).
   Der Inhalt von `u-boot.yaml` jenseits des Namens und die Version von `u-boot`, die das Gerüst erzeugte,
   werden nicht geprüft (Punkt 17).

6. **Schlüsselpaar.** `ssh/id_ed25519` und `ssh/id_ed25519.pub` unter `$PGR_DEV`. Zustand in `zustand` (nur lesen).
   - Beide fehlen: `schluessel` legt `ssh/` mit Modus 700 an und erzeugt das Paar ohne Passphrase mit dem
     Kommentar `<Projektname>`; die private Datei hat Modus 600. Besteht `ssh/` schon ohne die zwei Dateien, setzt `schluessel` seinen Modus auf 700 (Fall `p6-modus-ordner-vorhanden`, `ssh/` mit Modus 755 und leer; Mutation: das Setzen entfällt); andere Dateien in `ssh/` bleiben unberührt.
   - Beide vorhanden: sie werden genutzt, wenn (a) `ssh/` und die private Datei für Gruppe und andere weder lesbar
     noch schreibbar noch ausführbar sind, (b) `ssh-keygen -y -P ''` die private Datei ohne Passphrase liest und
     (c) dessen Ausgabe in Typ und Schlüsselteil der `.pub` entspricht. Sonst abgelehnt, Exit 4, die Meldung
     nennt den Pfad und welche der drei Bedingungen. Die Dateien bleiben in allen Fällen byteweise gleich und
     werden nie neu erzeugt, auch nicht bei Mängeln.
   - Nur eine der beiden Dateien: abgelehnt, Exit 4 (nie überschreiben).
   - `ssh` besteht und ist kein Verzeichnis (Datei, Symlink auf eine Datei): abgelehnt in `zustand`, Exit 4, die Meldung nennt den Pfad; bis dahin ist nichts geändert (Fall `p6-ssh-ist-datei`, `$PGR_DEV` danach unverändert, kein Aufruf von `u-boot`; Mutation: die Prüfung entfällt, dann scheitert `schluessel` erst nach `geruest`).
   Fälle `p6-*`: Paar fehlt (Rechte 700 und 600, Kommentar), Paar gültig unverändert, nur privat, nur öffentlich, Rechte
   von `ssh/` 750, Rechte der privaten Datei 640, Passphrase, `.pub` eines anderen Schlüssels, leere private
   Datei. Mutationen: überschreibt, prüft (a), (b), (c) nicht, setzt Rechte nicht, legt mit Passphrase an.

7. **`u-boot`.** In `geruest`, im Verzeichnis `$PGR_DEV`, in dieser Reihenfolge, mit diesen Argumenten:
   `init <Projektname> --devcontainer --sandbox --no-git`; `config set` mit
   `devcontainer.features.node.enabled true`, `devcontainer.features.sshd.enabled true`,
   `devcontainer.features.claude-code.enabled true`, `devcontainer.sandbox.nestedRuntime docker`,
   `devcontainer.sandbox.egress.enabled true`,
   `devcontainer.sandbox.egress.allow gcr.io,pkg-containers.githubusercontent.com`; `generate devcontainer`.
   Ein Aufruf mit Exit-Code 0 gilt als Erfolg, auch mit Warnungen auf stderr (Lockerungen des Profils, fehlender
   Remote). Ein Aufruf mit Exit-Code ungleich 0 beendet den Lauf mit Exit-Code 5 und der Meldung `geruest`;
   kein weiterer Aufruf folgt. Die Ausgabe von `u-boot` erscheint auf stdout und stderr unverändert und vollständig,
   das Ziel filtert sie nicht. `u-boot up` und `u-boot doctor` ruft das Ziel nie. Fälle `p7-*`: Aufrufe und Reihenfolge
   (Protokoll der Attrappe), ein Fall je Aufruf (acht), der scheitert (kein Folge-Aufruf, Exit 5, Kennzeichen der
   Ausgabe der Attrappe auf stderr), Warnung bei Exit 0 ist Erfolg, das Fortsetzen nach Punkt 5 ruft kein `init`.
   Mutationen: Argument weggelassen oder vertauscht, Fehler eines Aufrufs überlesen, Ausgabe verschluckt, `up` gerufen.
   Ein Fall gegen das echte `u-boot` (Gruppe Docker, ohne Container): `u-boot.yaml` trägt Namen, die drei Features,
   `nestedRuntime: docker` und die zwei Hosts; `.devcontainer/devcontainer.json` besteht.

8. **`mk-override.sh` und `override.json`.** Schritt `override`. Fehlt `mk-override.sh` im Verzeichnis, legt das Ziel es an
   (Modus 755); ein vorhandenes wird nicht überschrieben (Inhalt nicht geprüft, Punkt 17). Das Skript läuft im Verzeichnis
   mit `PGR_NAME`, `PGR_SSH_PORT`, `PGR_CPUS` und `PGR_MEM` in der Umgebung (Vorgaben eingesetzt). Das vom Ziel angelegte
   Skript schreibt `.devcontainer/override.json` als die von `u-boot` erzeugte `devcontainer.json` plus genau diese Zusätze:
   (a) `runArgs` um `--name=<PGR_NAME>`, (b) um `--cpus=<PGR_CPUS>`, (c) um `--memory=<PGR_MEM>`; (d) `appPort`
   gleich `["127.0.0.1:<PGR_SSH_PORT>:2222"]`; (e) `mounts` um einen Bind-Mount von `ssh/id_ed25519.pub` (unter
   `${localWorkspaceFolder}`) nach `/home/vscode/.ssh/authorized_keys` mit `readonly`; (f) `remoteEnv.CLAUDE_CODE_OAUTH_TOKEN`
   gleich dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}`. Alles andere der Datei bleibt gleich. **Nie ein Wert:** weder
   `override.json` noch eine Ausgabe des Ziels noch eine Datei unter `$PGR_DEV` enthält den Wert von `CLAUDE_CODE_OAUTH_TOKEN`
   oder `ANTHROPIC_API_KEY` aus der Umgebung; das Ziel liest, löst und übergibt keinen Zugang.
   **Umgebung der Kinder:** Das Skript entfernt `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` und `GIT_TOKEN` zu Beginn aus seiner eigenen Umgebung; kein Kindprozess (`u-boot`, `mk-override.sh`, `docker`, `devcontainer`, `ssh`) erbt sie. Der Verweis `${localEnv:…}` in `override.json` löst die `devcontainer`-CLI damit zu einem leeren Wert auf, und kein Zugang kann in die Ausgabe von `devcontainer up`, in `docker inspect` oder in die Prozessliste des Hosts gelangen (nicht gemessen und nicht nötig: der Wert ist nie vorhanden). Wer Claude Code im Container braucht, übergibt den Zugang in einem eigenen Aufruf (`devcontainer exec`), nicht über dieses Ziel.
   Fälle `p8-*` mit echtem `node` und einer minimalen `devcontainer.json`: je Zusatz (a) bis (f) einer, Rest gleich, Vorgaben,
   Wert-Sentinel in allen drei Variablen nirgends zu finden und in der Umgebung jeder Attrappe (Protokoll mit Umgebung) und von `mk-override.sh` nicht gesetzt (`p8-umgebung`), vorhandenes `mk-override.sh` bleibt. Mutationen: Zusatz fehlt oder
   falsch (Port, Adresse `0.0.0.0`, Mount ohne `readonly`), Verweis durch den Wert der Umgebung ersetzt, eine der drei Variablen wird nicht entfernt (je Variable eine Mutation), Rest der Datei verändert, Skript überschrieben.

9. **Container und zweiter Aufruf.** Schritt `zustand` liest, `start` handelt. Der Container „des Gerüsts“ ist der, dessen Label
   `devcontainer.local_folder` gleich `$PGR_DEV` ist (die Annahme, dass die CLI es setzt, misst der Implementer an einem Scratch-Container
   vor dem Code; trägt sie nicht, hält er an und gibt den Punkt dem Architect zurück).
   - Ein Container mit dem Namen `PGR_NAME`, der das Label nicht trägt: abgelehnt (Kollision), Exit 4; nichts wird angefasst.
   - Ein Container mit dem Label, dessen Name nicht `PGR_NAME` ist: abgelehnt, Exit 4, die Meldung nennt den vorhandenen
     Namen (zwei Container teilten die Volumes).
   - Ein Container mit Label und Name, dessen veröffentlichter Port 2222 nicht `127.0.0.1:<PGR_SSH_PORT>` ist: abgelehnt, Exit 4,
     die Meldung nennt den veröffentlichten Port. `PGR_CPUS` und `PGR_MEM` eines vorhandenen Containers werden nicht verglichen (Punkt 17).
   - Er läuft: kein `devcontainer up`, kein `override`; der Lauf geht zu `warten` und meldet Erfolg.
   - Er besteht, läuft nicht (`exited`, `created`): `override` und `devcontainer up --workspace-folder . --override-config .devcontainer/override.json`
     im Verzeichnis (Neustart; derselbe Container).
   - Es besteht keiner: ebenso (neuer Container; vorhandene Volumes bleiben, die Namen tragen die `devcontainerId` des Verzeichnisses).
   - **Randformen** (Entscheidung des Architect, Rückgabe des Implementers): (a) *Reihenfolge der Ablehnungen im Container:* Kollision des
     Namens (Name ohne Label) vor anderem Namen vor Port vor Zustand (b) und Verzeichnis (d); treffen Kollision und anderer Name zugleich zu,
     meldet das Ziel die Kollision. (b) *Zustand:* nur `running` gilt als läuft, nur `exited` und `created` als besteht, läuft nicht; jeder
     andere Zustand (`paused`, `restarting`, `removing`, `dead`) ist abgelehnt, Exit 4, die Meldung nennt Zustand und Namen, nichts wird
     angefasst (Grund: `devcontainer up` auf einen solchen Container ist nicht gemessen, und ein Neustart oder ein zweiter Container
     verschlimmerte den Zustand). (c) *Kein veröffentlichter Port 2222* gilt als abweichend: Exit 4, die Meldung nennt „keinen“.
     (d) *Container ohne Gerüst:* besteht ein Container mit dem Label und ist das Verzeichnis fehlend, leer oder ohne
     `.devcontainer/devcontainer.json`, ist der Lauf abgelehnt, Exit 4, die Meldung nennt das Verzeichnis; kein Schritt danach läuft
     (Grund: ein neu erzeugtes Gerüst und Schlüsselpaar passten nicht zu dem Container, der noch das alte hält). (e) *Docker-Fehler:*
     endet `docker ps` (nach Label oder nach Name) oder `docker inspect` in `zustand` mit Exit ungleich 0, ist der Lauf abgelehnt, Exit 4,
     die Meldung nennt den Aufruf; eine leere Ausgabe nach einem Fehler gilt nie als „kein Container“ (Grund: sonst startete das Ziel
     einen zweiten Container). Fälle an Attrappen: `p9-zustand-paused`, `-restarting`, `-removing`, `-dead` (je Zustand einer, kein `up`,
     Exit 4), `p9-kein-port` (Container ohne Bindung, Exit 4, Meldung mit „keinen“), `p9-ohne-geruest` (läuft, Verzeichnis leer; und
     gestoppt, Verzeichnis ohne `devcontainer.json`: je Exit 4, kein `u-boot`, kein `up`), `p9-docker-ps-label`, `p9-docker-ps-name`,
     `p9-docker-inspect` (Exit ungleich 0 und leere Ausgabe: Exit 4, kein `up`), `p9-doppelt` (Name ohne Label und anderer Container mit
     Label: Meldung der Kollision). Mutationen: ein Zustand wird als „läuft nicht“ behandelt (je Zustand, falls getrennt gelesen); die
     Prüfung des Ports lässt eine leere Bindung durch; die Prüfung (d) entfällt; der Fehler eines der drei Aufrufe wird als leer gelesen
     (je Aufruf); die Reihenfolge von Kollision und anderem Namen ist vertauscht.
   `devcontainer up` gilt als erfolgreich bei Exit-Code 0 **und** einer Ausgabezeile mit `"outcome":"success"`; sonst endet der Lauf
   mit Exit 5 (`start`), ohne Rückbau. Seine Ausgabe erscheint unverändert; er läuft ohne die drei Variablen aus Punkt 8 in der Umgebung (Fall `p9-umgebung`: die Attrappe von `devcontainer` sieht keine; Mutation: `env -u` entfällt). Ein Container, den Docker beim Fehler im Zustand `created`
   hinterlässt, ist beim nächsten Aufruf „besteht, läuft nicht“. **Stopp und Neustart (E6, Empfehlung):** kein zweites Ziel; Stopp ist
   `docker stop <Name>`, Neustart der zweite Aufruf, ebenso nach `docker rm` (neuer Container, Volumes bleiben; was im
   Dateisystem des Containers lag, ist weg). **Limits eines vorhandenen Containers** ändert das Ziel nie: kein Aufruf von `docker update`,
   auch nicht bei laufendem Container und abweichenden `PGR_CPUS` oder `PGR_MEM` (Entscheidung des Architect, Empfehlung des Planners;
   Grund: ein Ziel, das stillschweigend Limits eines laufenden Containers ändert, widerspräche „läuft: nichts tun“, und der Vergleich ist
   nach Punkt 17 (3) nicht zugesagt; der Weg ohne Neustart steht in Punkt 2). Fall `p9-kein-update`: in keinem Fall der Gruppe Attrappen
   enthält das Protokoll von `docker` einen Aufruf `update`, auch nicht bei laufendem Container mit `PGR_MEM`, das vom Container abweicht;
   Mutation: das Ziel ruft bei laufendem Container `docker update` mit den Werten. Das Ziel baut nach einer Änderung von `egress.allow` nicht um; der Nutzer ruft `u-boot generate devcontainer`
   und entfernt den Container von Hand. Attrappen-Fälle `p9-*` (Protokoll der Aufrufe): Kollision, Name weicht ab, Port weicht ab, läuft
   (kein `up`), gestoppt (`up`), `created` (`up`), keiner (`up`), `up` ohne `"outcome":"success"` trotz Exit 0, `up` Exit 1; Mutationen:
   Prüfung entfällt oder `up` bei laufendem Container. Fälle gegen reales Docker (Gruppe Docker): Erststart mit dem Gerüst aus Punkt 7
   (Container läuft, Name und Label stimmen), zweiter Aufruf bei laufendem Container (Startzeit des Containers unverändert),
   `docker stop` dann Aufruf (dieselbe Container-ID, läuft), `docker rm -fv` auf den Container ohne die Volumes und Aufruf (neue ID,
   dieselben Volume-Namen). Mutation des letzten: der Aufruf legt andere Volumes an (anderes Verzeichnis).

10. **Port belegt.** Läuft der Container des Gerüsts nicht (kein Container oder gestoppt), prüft `zustand` vor dem Start, ob eine
    Verbindung zu `127.0.0.1:<PGR_SSH_PORT>` zustande kommt (Verbindungsaufbau mit der Shell, ohne weiteres Werkzeug); kommt
    sie zustande, ist der Port belegt: abgelehnt, Exit 4, die Meldung nennt den Port. Ein Container, der den Port nicht veröffentlicht,
    sondern auf einer anderen Adresse oder ohne Listener belegt, wird so nicht erkannt; Docker meldet es dann in `start`, Exit 5
    (Punkt 17). Fall gegen reales Docker (Gruppe Docker): ein Container des Präfixes, der den Port veröffentlicht, und ein Aufruf mit
    diesem Port endet mit 4, und `docker ps` nennt keinen Container des Gerüsts; Mutation: die Prüfung entfällt (dann scheitert `start`
    mit Exit 5, der Fall wird rot) und die Prüfung läuft auch bei laufendem Container (Fall `p9-laeuft` wird rot).

11. **Bereitschaft.** Schritt `warten`. Bereit ist der Container, wenn `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf
    `127.0.0.1:<PGR_SSH_PORT>` den Befehl `true` mit Exit 0 ausführt. Der Aufruf lautet `ssh -i <Schlüssel> -o IdentitiesOnly=yes -o BatchMode=yes
    -o ConnectTimeout=5 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR -p <Port> vscode@127.0.0.1 true`.
    **A9 (Empfehlung des Architect, Entscheidung des Nutzers):** Die zwei Optionen für den Hostschlüssel gelten nur für diese Probe:
    Sie prüft die Erreichbarkeit der eigenen Schleife des Hosts, kein Vertrauen. Das Ziel schreibt und ändert nie `known_hosts`,
    `~/.ssh/config` oder `core.sshCommand` und setzt keine Einstellung für spätere Aufrufe; die Entscheidung über den Hostschlüssel
    bleibt beim Nutzer. Die Probe wiederholt sich jede Sekunde bis `PGR_WARTEN` Sekunden verstrichen sind; die Vorgabe 60 aus Punkt 2 wird hier beobachtbar (Fall `p11-vorgabe-warten` mit nicht gesetztem `PGR_WARTEN` und einer Attrappe von `ssh`, die nie gelingt: die Zahl ihrer Aufrufe liegt bei 60 oder 61; Mutation: Vorgabe 61 liefert 62 oder mehr; in Teil I ist der Wert nur in der Form geprüft, ein Mutant der Vorgabe bleibt dort grün, akzeptiertes Negativ bis Teil II). Erfolg: Meldung auf stdout
    `container-start: bereit <Name> 127.0.0.1:<Port>`, Exit 0. Kein Erfolg bis zur Zeitgrenze: Exit 5, die Meldung nennt `warten`
    und die letzte Ausgabe von `ssh`; der Container bleibt stehen. Ob `dockerd` im Container bereit ist, prüft das Ziel nicht
    (`devcontainer up` endet nach `postStartCommand`; Punkt 17). Fälle `p11-*` an einer Attrappe von `ssh`: Erfolg im ersten Versuch,
    Erfolg nach drei Versuchen, nie (mit `PGR_WARTEN=2`, Exit 5, letzte Ausgabe in der Meldung, Dauer nicht unter 2 Sekunden), je Option
    des Aufrufs eine (Protokoll der Attrappe: Schlüssel, `IdentitiesOnly`, `BatchMode`, Port, Benutzer, Adresse, `UserKnownHostsFile=/dev/null`),
    `known_hosts` und `~/.ssh/` des Hosts nach dem Lauf unverändert (Mutation: ohne `UserKnownHostsFile=/dev/null`). Fall gegen reales Docker:
    die Probe meldet bereit mit dem Schlüssel des Gerüsts und scheitert (Exit 5) mit einem anderen.

12. **Abbruch (Teil II).** Zugesagt erst mit Teil II; im Zwischenstand (Punkt 13) ist das Verhalten bei `SIGINT` und `SIGTERM` nicht zugesagt und hat keinen Fall (akzeptiertes Negativ; Grund: der Zwischenstand endet nach `override` ohne langlaufendes Kind, und die Falle gehört zum Start). Der Fall im Schritt `geruest` und die Wiederholung gehören zu Teil II. `SIGINT` und `SIGTERM` beenden den Lauf: Das Skript beendet seinen laufenden Kindprozess (`u-boot`, `devcontainer`, `ssh`),
    meldet auf stderr `container-start: <Schritt>: abgebrochen` mit dem Schritt, in dem es war, und endet mit 130 beziehungsweise 143. Es
    baut nichts zurück. Jeder Schritt ist wiederholbar (Punkte 5, 6, 9), mit zwei akzeptierten Negativen: (1) Bricht `u-boot init` ab und hinterlässt
    Dateien ohne `u-boot.yaml`, ist das Verzeichnis beim nächsten Aufruf fremd (Punkt 5) und der Nutzer räumt es; (2) bricht der Lauf in
    der Lücke zwischen privater und öffentlicher Datei von `ssh-keygen` ab, meldet der nächste Aufruf „nur eine Datei“ (Punkt 6); Grund beider:
    ein Zeitfenster von Millisekunden bis Sekunden in Werkzeugen, die nicht das Ziel sind, und ein Rückbau wäre ein Mechanismus für einen
    Fall, den die Meldung dem Nutzer nennt. `SIGKILL` und Ausfall des Hosts fängt nichts auf. Fälle `p12-sigint`, `p12-sigterm` an einer
    Attrappe von `devcontainer`, die wartet (Exit-Code, Meldung mit Schritt `start`, der Kindprozess lebt nicht mehr, Verzeichnis
    unverändert gegenüber dem Stand vor `start`), dazu ein Fall im Schritt `geruest` (Attrappe von `u-boot`) und die Wiederholung
    danach (`p12-wiederholung`: der nächste Aufruf setzt fort, Punkt 5). Mutationen: die Falle beendet das Kind nicht, die Falle
    entfällt für `INT` oder `TERM` (äquivalent, wenn bash allein 130 beziehungsweise 143 ergibt: die Gegenprobe führt es dann
    als Mutant ohne roten Fall, wie `SPEC-057` Punkt 7; gemessen vom Implementer).

13. **Reihenfolge der Gründe und Ausgang.** Treffen mehrere Gründe zu, gilt diese Reihenfolge, der erste beendet den Lauf (innerhalb
    von `eingaben` und `werkzeuge` werden alle Verstöße des Schrittes gemeldet, Punkte 2 und 3):
    1. `eingaben` (Exit 2), 2. `werkzeuge` (3), 3. in `zustand` (4): Lage des Pfades im Repo, Zustand des Verzeichnisses (Punkt 5), Zustand
    des Schlüsselpaars (6), Zustand des Containers (9), Port belegt (10), 4. Fehler eines Schrittes `geruest`, `schluessel`, `override`,
    `start`, `warten` (5), 5. Signal (130, 143). **Signal und Fehler eines Schrittes (5) sind Ereignisse, keine Zustände:** das erste beendet den Lauf. Trifft ein Signal ein, während
    ein Kindprozess läuft, beendet der Abbruch das Kind und der Lauf endet mit 130 oder 143, nie mit 5 (Fall `p12-*` des Abbruchs, Punkt 12);
    hat ein Schritt schon mit 5 geendet, bleibt es bei 5. Nicht zugesagt und nicht prüfbar ist das Ergebnis bei Gleichzeitigkeit beider
    (akzeptiertes Negativ; Grund: kein deterministischer Fall); geprüft wird nur die Reihenfolge der Ereignisse (Fall `p13-fehler-vor-signal`: Schritt scheitert, danach
    kein Signal, Exit 5; Mutation: der Fehler wird als Signal gemeldet). Exit-Code 0 heißt „bereit“ (auch beim zweiten Aufruf); 1 und alle anderen sind nicht
    belegt. **Zwischenstand, nicht zugesagt:** Solange Teil II nicht geliefert ist, endet ein Lauf, der alle Schritte bis `override` besteht, mit Exit-Code 6 und der Zeile `container-start: override: Teil II fehlt, Gerüst liegt, kein Container gestartet` auf stderr; er schreibt nie `bereit` und endet nie mit 0. Mit Teil II entfallen Exit 6 und die Zeile (kein Code und keine Zeile des Vertrags nennen ihn danach); die Fälle des Zwischenstandes (`p13-zwischenstand`, Mutation: Exit 0 oder Zeile `bereit`) fallen mit ihm weg. Fälle `p13-*`: je Nachbarpaar der Reihenfolge ein Fall, in dem beide Gründe zutreffen und der frühere gemeldet wird (Eingabe vor
    Werkzeug; Werkzeug vor Pfad im Repo; Pfad im Repo vor fremdem Verzeichnis; fremdes Verzeichnis vor Schlüsselmangel; Schlüsselmangel vor
    Kollision des Containers; Kollision vor belegtem Port; Kollision vor anderem Namen, Punkt 9 (a)) und der Exit-Code jedes Grundes; Mutation: die beiden Prüfungen vertauscht.

14. **Aufrufarten des Treibers.** Das Ergebnis ist der Exit-Code des Skripts und die Meldung `bereit`; kein Stand und kein Hook glaubt es.
    Das Rezept von `container-start` ist eine Zeile ohne `$(MAKE)` und ohne `+`, das Ziel `.PHONY` ohne Voraussetzungen. Unter `-n`
    zeigt `make` das Rezept und führt es nicht aus: kein Verzeichnis, kein `u-boot`, kein `ssh-keygen`, kein `docker`, kein `devcontainer`.
    Unter `-q` und `-t` ebenso nichts (Exit-Code von `-q` ist der von `make`, nicht ausgewertet) und kein Zeitstempel, weil das Ziel
    keine Datei ist. Unter `-i` überliest `make` den Exit-Code des Rezepts: Der Exit-Code von `make` belegt dann nichts, und die
    Zusage ist die Meldung: `container-start: bereit` steht nur nach einem Erfolg, bei jedem Fehler stattdessen die Zeile des
    Schrittes auf stderr. `MAKEFLAGS` aus der Umgebung gilt wie das Flag, auch in Kombination; bei `-n`, `-q` oder `-t` zusammen mit `-i` gilt
    `-n`, `-q`, `-t` (nichts läuft). `-k` und `-j` ändern an diesem einen Ziel nichts: gleiches Ergebnis wie ohne. Gleichzeitige Aufrufe
    mit demselben `PGR_DEV` sind nicht zugesagt (Punkt 17). Der Erfolgsfall eines Laufs ist in Teil I der Zwischenstand aus Punkt 13 (Exit 6 von `make` als Fehler des Rezepts, Zeile des Zwischenstandes, keine Zeile `bereit`), mit Teil II der Lauf mit `bereit`. Je Flag `-n`, `-q`, `-t`, `-i`, `-k`, `-j`, je `MAKEFLAGS` aus der Umgebung
    mit demselben Flag und für `-n` zusammen mit `-i` ein Fall `p14-*` an Attrappen (Protokoll der Aufrufe leer bei `-n`, `-q`, `-t`; bei
    `-i` ein Fehlerfall ohne Zeile `bereit`; bei `-k` und `-j` ein Erfolgsfall im eben genannten Sinn); Mutation: ein Rezept mit `$(MAKE)` oder `+`.

15. **Gegenprobe.** `container-start-gegenprobe` hat zwei Gruppen. *Gruppe Attrappen* (alle Punkte außer den Fällen gegen reales Docker)
    läuft in Temp-Verzeichnissen mit Attrappen von `docker`, `u-boot`, `devcontainer` und `ssh` in einem vorangestellten Pfad, ohne Docker
    und ohne Netz, mit echtem `node` und `ssh-keygen`; jede Attrappe protokolliert Aufruf und Argumente. *Gruppe Docker* (die Fälle gegen
    reales Docker aus den Punkten 7 bis 11) fährt gegen den Daemon des Hosts und höchstens einen Container des Gerüsts zugleich. Regeln der
    Gruppe Docker (Zuordnung: (a) für Verzeichnisse und (f) für Verzeichnisse und Temp-Pfade gelten schon in Teil I, mit dem Fall gegen das echte `u-boot` aus Punkt 7, der eine Attrappe von `docker` nutzt, weil Schritt `werkzeuge` nur `info` und `version` ruft; (a) für Container und Image, (b), (c), (d), (e) und (f) für Container und Volumes erst mit Teil II; (e) hat in Teil I keinen Fall, weil dort kein Name von außen kommt: das Verzeichnis stammt aus `mktemp` mit dem Präfix, der Vergleich vor dem Entfernen ist Absicherung ohne Fall (akzeptiertes Negativ; Grund: ein Fall bräuchte eine Eingabe, die es in Teil I nicht gibt)): (a) Namenspräfix `pgr-cs-` für jeden Container, jedes Verzeichnis (unter `mktemp -d`), jedes Image (`vsc-pgr-cs-`);
    (b) der Port ist einer, der zum Zeitpunkt frei ist, nie 2222 und nie ein Port eines Containers ohne Präfix; (c) `PGR_CPUS=2` und
    `PGR_MEM=2g`; (d) Entfernen nur per Name: `docker rm -fv`, danach die Prüfung, dass keines der zuvor aus `docker inspect` gelesenen Volumes
    (Typ `volume`) bleibt (`SPEC-057` Punkte 3 und 4; ein Rest wird gemeldet, per Name entfernt und färbt die Gegenprobe rot), `docker rmi`
    nur für Images mit dem Präfix; kein `prune`; (e) kein Name ohne Präfix wird angefasst (Fall: die Gegenprobe bricht ab, wenn ein Name aus
    der Umgebung das Präfix nicht trägt); (f) das Aufräumen läuft bei Ende, `SIGINT` und `SIGTERM`; (g) entfernt wird nur,
    was der eigene Lauf angelegt hat: die Gegenprobe merkt sich die Namen ihrer Verzeichnisse, Container, Volumes und Images und sucht nie
    per Präfix; gleichzeitige Läufe am Host mit demselben Präfix `pgr-cs-` beeinflussen sich nicht (Entscheidung des Architect nach Befund
    V-3 des Verifiers; Fall `p15-fremdes-verzeichnis`: ein fremdes leeres Verzeichnis `pgr-cs-*` bleibt stehen, auch nach `SIGINT` und
    `SIGTERM` des eigenen Laufs; Mutation: Aufräumen per Präfix statt per gemerktem Namen). Eine Vorbedingung „kein `make gates` läuft“
    gibt es nicht: gleichzeitige Läufe auf einem Daemon sind nicht zugesagt (`SPEC-052`
    Punkt 8); die Regel „nie gleichzeitig mit dem Lauf eines anderen Slice und nie gegen den Container, in dem gearbeitet wird“ richtet sich an
    den Aufrufer. Die Mutanten ersetzen genau eine Stelle in einer Kopie des Skripts; trifft die Ersetzung nichts oder mehr als eine Stelle,
    bricht die Gegenprobe ab. **Voraussetzung am Host:** Die Gegenprobe braucht `u-boot` (echt, Fall aus Punkt 7), `node` und `ssh-keygen` im Pfad; in Teil II zusätzlich Docker und `devcontainer`. Fehlt eines, endet sie mit Exit-Code ungleich 0 und nennt das Werkzeug, bevor ein Fall läuft; sie überspringt keinen Fall (R1: ein Lauf, der ohne Voraussetzung grün würde, belegte den Fall nicht). Die Gegenprobe ist ein Ziel für den Host (Punkt 16); in einem Dev-Container läuft sie nicht, und es gibt dort keine Erkennung und kein Überspringen. Fall `p15-voraussetzung` je Werkzeug mit einem Pfad ohne dieses; Mutation: der Fall wird bei fehlendem Werkzeug übersprungen.
    Die Laufzeit ist zu messen und im Plan zu nennen.

16. **Host-Ziel, kein Gate.** `container-start-gegenprobe` ist `.PHONY` und ein reines Ziel für den Host (Entscheidung des Nutzers vom 2026-10-18: die Gegenproben werden aufgeteilt, die mit `u-boot` und `devcontainer` laufen nur am Host). Sie hängt nicht an `GATE_CHECKS`, steht nicht in der Ausnahmeliste der schnellen Kette (`SPEC-052` Punkte 2 und 4) und gehört nicht zum Sammelziel der Gegenproben, die im Dev-Container laufen; kein Hook und kein Stand glaubt ihr Ergebnis. Sie steht in der Liste für den Host und gehört zum Sammelziel `gegenproben-host` ([`SPEC-060`](SPEC-060-gegenproben.md) Punkte 1 und 4). Dass sie läuft und grün ist, belegt der Bericht des Slice, der Teil I oder Teil II ändert: Lauf am Host mit Mutationsprotokoll (Regel an den Aufrufer, kein Fall; Grund: ein Gate, das sie erzwänge, wäre das Gate, das diese Entscheidung verlässt). Sie legt ihre Verzeichnisse in Temp-Pfaden an. `container-start` selbst steht nicht in `GATE_CHECKS`. Fall `p16-keine-kette` in einem Temp-Baum mit Makefile und Fragmenten: `container-start-gegenprobe` hängt nicht an `GATE_CHECKS`, steht nicht in der Ausnahmeliste, `container-start` steht nicht in `GATE_CHECKS`, `make` endet bei rotem Skript mit Exit-Code ungleich 0; Mutationen: sie hängt an `GATE_CHECKS`, sie steht in der Ausnahmeliste.

17. **Grenzen.** (1) Nicht geprüft: Version von `devcontainer`, `node` und `ssh`; die Herkunft von Basis-Image, Features und APT-Repository;
    ein Build ohne warmen Cache des Hosts. (2) Inhalt von `u-boot.yaml` jenseits des Namens, Version von `u-boot`, die das Gerüst erzeugte, und Inhalt
    eines vorhandenen `mk-override.sh`: das Ziel startet das Gerüst, wie es ist. (3) `PGR_CPUS` und `PGR_MEM`: kein Wert wird geprüft, einer
    eines vorhandenen Containers nicht verglichen und nicht angepasst (Punkt 9); die Vorgaben sind nur gemessen, wie Punkt 2 sagt; ein ungültiger Wert endet als Fehler von Docker in `start` (Exit 5, ein Container im Zustand
    `created` bleibt). (4) Ob `dockerd` im Container bereit ist, und das, was im Container liegt. (5) Ein Port, der auf einer anderen Adresse
    als `127.0.0.1` oder ohne Listener belegt ist (Punkt 10). (6) Gleichzeitige Aufrufe mit demselben `PGR_DEV` (Grund: ein Nutzer, ein Terminal;
    eine Sperre brauchte ein Aufräumen nach `SIGKILL`). (7) Ein Hostschlüssel, der sich nach einem Neubau ändert, ist Sache des Nutzers
    (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`). (8) `nestedRuntime: docker` ist festgeschrieben; die Wahl zwischen `docker` und
    `podman` bleibt offen (`BEO-REPO/nested-runtime-podman-oder-docker-offen`). (9) Wiederholbarkeit ist an den Punkten 5, 6, 9 und 12 gebunden;
    ein Verzeichnis, das jemand von Hand verändert, ist nicht zugesagt. (10) Ob das Dokument `docs/developer/devcontainer.md` die Schritte
    richtig nennt, prüft kein Gate.
