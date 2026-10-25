# SPEC-069 — Einrichtung im laufenden Dev-Container (`container-einrichten`, `container-einrichten-gegenprobe`)

`make container-einrichten` richtet am Host für einen **laufenden** Dev-Container in einem Aufruf ein, was
die Arbeit an einem Slice im Container braucht: das Bare-Repo, den eingelieferten `main`, den Arbeitsklon,
dessen Identität, den Commit-Träger und den Zweig des Slice. `make container-einrichten-gegenprobe` prüft
das Werkzeug. Das Werkzeug ist kein Gate: Es steht im Abschnitt *Werkzeuge* von
[`harness/README.md`](../README.md), nicht in `GATE_CHECKS`, und kein Hook und kein Stand glaubt sein
Ergebnis. Diese Kennung legt fest, was es zusagt, woran es scheitert und was die Gegenprobe belegt.

**Stand.** Die Entscheidungen dieser Kennung gelten:
[ADR-0057](../../docs/plan/adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) ist
angenommen (Entscheidung des Nutzers vom 2026-10-21). `ssh` als Voraussetzung am Host für dieses Ziel
ist Entscheidung des Nutzers vom 2026-10-17 („ja“), die Obergrenze für jede Gegenprobe seine
Entscheidung vom 2026-10-21 (600 s; seit 2026-10-25 1200 s, Punkt 17)
([ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)).

**Was es nicht ist.** Es legt kein Gerüst an, erzeugt kein Schlüsselpaar, startet, stoppt und entfernt
keinen Container und prüft keine Voraussetzung von `u-boot`, der `devcontainer`-CLI oder `node`: Das sagt
[`SPEC-058`](SPEC-058-container-start.md) zu, und diese Kennung setzt dessen Punkte 2 (Eingaben, Vorgaben)
und 11 (ein Container, der per `ssh` als `vscode` antwortet) voraus. Es löst, liest und übergibt keinen
Zugang von Claude Code (Punkt 13). Es ruft `docker` nie auf (Punkt 1). Ob `dockerd` im Container bereit
ist, sagt weder `SPEC-058` Punkt 17 (4) noch diese Kennung (Punkt 20 (4)).

Der Vertrag hat zwei Teile, damit ein Schnitt an einer Naht liegt: **Teil 1** (Punkte 2 bis 5) prüft
Eingaben, Werkzeuge und Zustand, ohne dass sich am Host oder im Container etwas ändert; **Teil 2**
(Punkte 6 bis 13) handelt. Punkte 14 bis 20 gelten für beide. Liefert ein Slice nur Teil 1, trägt er den
Zwischenstand nach dem Muster von [`SPEC-058`](SPEC-058-container-start.md) Punkt 13 nach (eigener
Exit-Code, eigene Zeile, keine Erfolgsmeldung); solange ein Slice beide Teile liefert, nennt diese Kennung
keinen Zwischenstand.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element
einer Aufzählung gibt es einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die
Gegenprobe verlangt genau diese Menge roter Fälle (mehr und weniger machen sie rot).

**Plattform.** Zugesagt und geprüft ist ein **Linux-Host** mit bash ≥ 4.4 und GNU-Werkzeugen, OpenSSH und
`git`, gegen einen Linux-Container des Gerüsts aus `SPEC-058`. macOS- und Windows-Hosts sind **nicht
zugesagt** (Grenze, `AGENTS.md` §3.12 *Die Plattform gehört dazu*; die drei Betriebssysteme von
[`LH-QA-03`](../../spec/lastenheft.md#lh-qa-03--portabilität) gelten dem Produkt, nicht diesem Werkzeug).
Plattformabhängig sind hier: die Optionsformen von `ssh` (OpenSSH; `ssh.exe` unter Windows ist nicht
gemessen), der Transport der Identität über den Standardeingang (ein Zeilenende `CRLF` liefert einen Wert
mit `CR`; deshalb lehnt Punkt 5 Steuerzeichen ab), die Form des Pfades des Host-Repos und die
GNU-Optionen der Werkzeuge, die das Skript ruft. Das Ziel lehnt einen fremden Host nicht ab, damit ein
Nutzer es dort auf eigene Gefahr versuchen kann; was dort geschieht, ist nicht gemessen.

**Warum 4.4** (Entscheidung des Architect, 2026-10-21): Das Werkzeug wartet auf seinen Kindprozess über `$!`
nach einer Prozess-Substitution; dass `$!` sie benennt und `wait` sie annimmt, führt bash ab 4.4.
**Akzeptiertes Negativ, mit Grund:** Unter bash 4.0 bis 4.3 endet der Lauf im Lesezug des Schrittes
`zustand` mit einer Meldung von bash und einem Exit-Code, den Punkt 14 nicht belegt — also **vor** der
ersten Änderung an Host und Container (Punkt 4) und sichtbar für den, der den Aufruf abgesetzt hat. Das Ziel
prüft seine bash-Fassung **nicht selbst**: `BASH_VERSINFO` ist nicht setzbar, ein Fall bräuchte deshalb eine
Ersatz-Variable wie `GP_BASH_VERSION` in
[`SPEC-061`](SPEC-061-gegenproben-nach-einlieferung.md) und damit eine sechste Eingabe mit eigenen
Randformen (Punkt 2) — in einer Gegenprobe ohne Nutzer-Eingaben trägt das, in einem Ziel nicht. Erscheint
ein Host mit einer älteren Fassung, bekommt der Schritt `werkzeuge` die Prüfung mit Exit 3 und einen Fall,
der die Fassung ersetzt.

**Randform → Punkt.** Die Nummern der Randformen sind die des Slice-Plans
`slice-harness-devcontainer-einrichten` §6.

| Randform des Plans | Punkt |
|---|---|
| 1 Bare-Repo besteht schon | 5 (Zustand), 6 |
| 2 Klon besteht schon | 5 (Zustand), 8 |
| 3 Remote besteht schon | 5 (Zustand), 7 |
| 4 Container nicht erreichbar | 5 (Lesezug), 14 |
| 5 Identität am Host leer | 5 (Identität), 14 |
| 6 Identität mit Sonderzeichen | 5 (Identität), 9 |
| 7 Identität nur lokal, nie global | 9 |
| 8 anderer Port, mehrere Instanzen | 2, 3 (Projektname), 7 (URL) |
| 9 Hostschlüssel unbekannt oder geändert | 5 (Lesezug), 20 (5) |
| 10 Netz fehlt bei `make traeger-fetch` | 12 |
| 11 Abbruch mitten im Ablauf, Wiederholung | 13, 15 |
| 12 `main` am Host fehlt | 5 (Host-Repo) |
| 13 Push abgewiesen | 7 |
| 14 Reihenfolge der Gründe, Exit-Codes | 14 |
| 16 Zweig je Slice im Klon | 2, 11 |
| 25 Stelle des Zugangs im Ablauf | 13 |
| Aufrufarten des Treibers | 16 |
| Plattform | Kopf dieser Kennung, 19 |
| Zeit und Platz der Gegenprobe | 17, [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16 |
| Start der Gegenprobe mit ignoriertem `SIGINT` | 17 (Start), 19, [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 17 |

1. **Gegenstand und Ort.** Das Ziel `container-einrichten` ist im Repo (Fragment
   `harness/mk/container-einrichten.mk`, Skript `tools/harness/container-einrichten.sh`), `.PHONY`, kein
   Gate. Es schreibt am Host **allein** das Remote des Host-Repos (Punkt 7) und im Container **allein**
   unter `/workspaces/<Projektname>/` (Punkte 6 bis 12). Es ändert nie: `core.sshCommand` oder eine andere
   Einstellung des Host-Repos außer dem Remote, `~/.ssh`, eine `known_hosts`-Datei, eine Datei im
   Arbeitsbaum des Host-Repos, die Git-Konfiguration auf Ebene `--global` oder `--system` am Host oder im
   Container, und `$PGR_DEV` samt Gerüst und Schlüsselpaar. Es ruft nie auf: `docker`, `u-boot`,
   `devcontainer`, `git push --force` (auch nicht `--force-with-lease`, `--mirror`, `--delete`),
   `git reset`, `git checkout -f`, `git clean`, `git merge`, `git pull`, `git rebase`, `git remote
   remove`, `git remote set-url`, `git branch -D`, `rm -rf`. Die Gegenprobe zeigt es an Attrappen, die
   jeden Aufruf protokollieren (Fall `p1-keine-zerstoerung` je Element der zwei Aufzählungen; Mutation: ein
   Schritt ruft das genannte Kommando).

2. **Eingaben.** Sie kommen aus der Umgebung. Ein Wert, der gesetzt und leer ist, ist ungültig (kein stilles
   Einsetzen der Vorgabe); nicht gesetzt heißt Vorgabe. Eine Ausnahme nennt die Tabelle.

   | Variable | Vorgabe | gültig |
   |---|---|---|
   | `PGR_DEV` | keine, Pflicht | absoluter Pfad (beginnt mit `/`) **aus `A-Za-z0-9._+-` und `/`** — dieselbe Zeichenmenge wie [`SPEC-058`](SPEC-058-container-start.md) Punkt 2, hier erneut geprüft (Grund unter der Tabelle) —, vorhanden und ein Verzeichnis; die übrigen Teile der Form (Segment `.`, `..`, leer, `/` am Ende) prüft `SPEC-058` Punkt 2 beim Anlegen des Gerüsts, und der Basisname fällt hier mit dem Projektnamen aus Punkt 3 zusammen |
   | `PGR_SSH_PORT` | `2222` | ganze Zahl 1 bis 65535, Dezimalziffern, keine führende Null, kein Vorzeichen (wie `SPEC-058` Punkt 2) |
   | `PGR_REMOTE` | `container` | aus `A-Za-z0-9._-`, erstes Zeichen `A-Za-z0-9`, höchstens 64 Zeichen |
   | `PGR_ZWEIG` | keine, Pflicht | von `git check-ref-format --branch` angenommen und nicht `main` (Grund: der Klon arbeitet nie auf `main`, Dokument §9); fehlt `git` am Host, bleibt **allein diese eine** Prüfung aus (Ausnahme, in Punkt 14 begründet), die übrigen Formen prüft `eingaben` auch dann |
   | `PGR_SSH_OPTS` | leer | zusätzliche Optionen für jeden `ssh`-Aufruf des Ziels, an Leerraum geteilt und als einzelne Argumente übergeben, nie von einer Shell ausgewertet (**das gilt den `ssh`-Aufrufen dieses Ziels; im Push reist der Wert nach Punkt 7 in `GIT_SSH_COMMAND`, Grenze in Punkt 20 (7)**); **gesetzt und leer ist gültig** und gleich nicht gesetzt (Grund: eine leere Optionsliste ist die Vorgabe); ein Zeilenumbruch oder ein anderes Steuerzeichen im Wert ist ungültig |

   **Warum `PGR_DEV` hier erneut geprüft wird** (Entscheidung des Architect, 2026-10-21): Der Wert kommt aus
   der Umgebung **dieses** Aufrufs; dass `container-start` ihn je gesehen hat, ist nicht gesagt, und die
   Eingabeprüfung eines anderen Ziels ist eine Aussage über einen Lauf, der nicht stattgefunden haben muss.
   Der Pfad reist in Punkt 7 in `GIT_SSH_COMMAND`, und `git` gibt diese Zeichenkette an eine Shell
   (Punkt 20 (7)); ein Zeichen außerhalb der Menge würde dort ausgewertet, und zwar im Schritt, der
   schreibt. Die Menge ist die des Gebers und schließt damit keinen Pfad aus, unter dem ein Gerüst entstehen
   konnte.

   Schritt `eingaben` endet bei einem Verstoß mit **Exit-Code 2**, ändert nichts und ruft kein Werkzeug auf.
   Er meldet **alle** Verstöße, je einen in einer Zeile, die die Variable nennt, in der Reihenfolge der
   Tabelle. Je Form ein Fall und eine Mutation, die nur ihre Prüfung entfallen lässt: `PGR_DEV` nicht
   gesetzt, leer, relativ, ein Zeichen außerhalb der Menge (`;`), ein Leerzeichen, fehlend, keine
   Datei-Art Verzeichnis; `PGR_SSH_PORT` leer, `abc`, `02222`,
   `+22`, `0`, `65536`, Grenzen `1` und `65535` angenommen; `PGR_REMOTE` leer, mit `/`, mit Leerzeichen,
   erstes Zeichen `-`, 65 Zeichen; `PGR_ZWEIG` nicht gesetzt, leer, `main`, `a..b`, `-x`, `x y` und ein Name, der mit
   `refs/` beginnt (abgelehnt, auch wenn `git check-ref-format --branch` ihn annimmt: eigener Fall, Grund:
   ein Zweigname mit `refs/` benennt im Klon nicht den Zweig, den der Aufrufer meint), Grenzen `a` und
   `slice-harness-devcontainer-einrichten` angenommen; `PGR_SSH_OPTS` nicht gesetzt, leer, eine Option, zwei Optionen, mit Tabulator als Trenner,
   mit Zeilenumbruch abgelehnt. Nicht gesetzte Variablen ergeben die Vorgaben (Fall `p2-vorgaben`, sichtbar
   an der URL des Remotes und am `ssh`-Aufruf; Mutation: andere Vorgabe). Fälle `p2-*`.

3. **Werkzeuge und Projektname.** Schritt `werkzeuge`: `ssh` und `git` müssen am Host im Pfad liegen; ein
   Verstoß endet mit **Exit-Code 3** und ändert nichts. Der Schritt meldet beide Verstöße, je einen in einer
   Zeile mit dem Namen des Werkzeugs, in der Reihenfolge `ssh`, `git`. Versionen prüft das Ziel nicht
   (Punkt 20 (1)).

   **Warum zwei Werkzeuge am Host:** Dieses Ziel ist kein Gate, steht nicht in `GATE_CHECKS` und richtet
   einen Dev-Container **von außen** ein; es kann das nicht in einem Container tun. `ssh` als Voraussetzung
   am Host ist die Entscheidung des Nutzers vom 2026-10-17 (Kopf dieser Kennung, *Stand*), `git` liegt am
   Host, weil der Klon dort liegt. Dieselbe Lage trägt `container-start` mit `docker`, `u-boot`,
   `devcontainer`, `node`, `ssh` und `ssh-keygen` ([`SPEC-058`](SPEC-058-container-start.md) Punkt 3,
   [ADR-0045](../../docs/plan/adr/0045-container-start-als-werkzeug-im-repo.md)). Die Host-Werkzeuge je Ziel
   stehen in dessen Vertrag und in der Zeile des Ziels in [`harness/README.md`](../README.md), nicht im
   Briefing; dass `AGENTS.md` §3.1 sie nicht führt, ist bekannt und von diesem Vertrag nicht geändert.

   Fälle `p3-*`: je Werkzeug fehlt, beide fehlen (Reihenfolge), `git` fehlt und
   `PGR_ZWEIG` ist allein nach `check-ref-format` ungültig (`p3-fehlt-git-zweig-ungueltig`: Exit 3, die
   Ausnahme aus Punkt 14); Mutationen: das Werkzeug wird nicht geprüft, Reihenfolge vertauscht, der
   Vorbehalt `command -v git` in `eingaben` entfällt (dann meldet `eingaben` auch ohne `git` den
   Zweignamen, und die drei Fälle mit fehlendem `git` werden rot).

   **Der Projektname** kommt aus `$PGR_DEV/u-boot.yaml`, aus der ersten Zeile der Form `    name: <Name>`
   (vier Leerzeichen, wie `SPEC-058` Punkt 5 sie verlangt). Er muss aus `a-z0-9-` bestehen, mit `a-z0-9`
   beginnen und gleich dem Basisnamen von `PGR_DEV` sein; sonst ist der Lauf abgelehnt (Exit 4, die Meldung
   nennt Datei und gefundenen Wert). Fehlt die Datei oder die Zeile, ist der Lauf ebenso abgelehnt mit der
   Meldung, dass kein Gerüst vorliegt, und dem Verweis auf `make container-start`. **Grund:** Der Name
   bestimmt den Pfad im Container (`/workspaces/<Projektname>/`); aus der Datei gelesen statt aus dem Wert
   der Variablen gebildet, trägt ihn die Quelle, die `u-boot` geschrieben hat (`SPEC-058` Punkt 5), und ein
   Pfad kann nicht aus einem Wert entstehen, den niemand geprüft hat. Fälle `p3-name-*`: Datei fehlt, Zeile
   fehlt, Name mit Großbuchstabe, Name mit Unterstrich, Name ungleich Basisname, zwei `name:`-Zeilen (die
   erste gilt), Zeile mit anderem Einzug (gilt nicht); Mutationen: Name aus dem Basisnamen gebildet, Vergleich
   entfällt, Einzug nicht geprüft.

4. **Ablauf.** Die Schritte laufen in dieser Reihenfolge, der erste Verstoß beendet den Lauf: `eingaben`
   (2), `werkzeuge` (3), `zustand` (5), `bare` (6), `einliefern` (7), `klon` (8), `identitaet` (9), `hooks`
   (10), `zweig` (11), `traeger` (12). **Bis einschließlich `zustand` ändert das Ziel nichts:** weder am
   Host noch im Container, und es ruft von `git` nur lesende Kommandos (`rev-parse`, `config --get`,
   `remote get-url`, `check-ref-format`) und von `ssh` nur den einen Lesezug aus Punkt 5. Fall
   `p4-nichts-geaendert` je Verstoß der Punkte 2, 3 und 5: Host-Repo (Konfiguration, Remotes, Arbeitsbaum),
   `$PGR_DEV` und das Protokoll der Attrappen nennen keinen anderen Aufruf; Mutation: ein Schritt schreibt
   vor seiner Prüfung. Jede Fehlermeldung nennt den Schritt (`container-einrichten: <Schritt>: <Grund>`,
   auf stderr); auf stdout steht von der eigenen Ausgabe des Ziels nur die Schlussmeldung des Erfolgs
   (Punkt 13). Die Ausgabe eines gerufenen Werkzeugs (`ssh`, `git`, `make` im Container) erscheint
   unverändert auf dem Strom, auf dem das Werkzeug sie schrieb, und ist davon ausgenommen.

5. **Zustand — ein Lesezug, alles vor der ersten Änderung.** Schritt `zustand` liest zuerst am Host, dann
   mit **einem** `ssh`-Aufruf im Container. Jeder Grund endet den Lauf mit **Exit-Code 4**; die Reihenfolge
   ist die dieser Liste (Punkt 14).

   1. **Gerüst** (lokal): Projektname nach Punkt 3; `$PGR_DEV/ssh/id_ed25519` vorhanden und lesbar.
   2. **Host-Repo** (lokal): Das Arbeitsverzeichnis liegt in einem Git-Repo (`git rev-parse
      --show-toplevel`); der Ref `refs/heads/main` existiert. Ein `detached HEAD` und ein anderer
      ausgechecker Zweig sind kein Grund: Das Ziel liefert immer `refs/heads/main` ein und hängt nicht am
      ausgecheckten Zweig (Randform 12).
   3. **Identität** (lokal): `git config --get user.name` und `git config --get user.email` im Host-Repo
      — **alle Ebenen**, auch `--global` und `--system` (Randform 7: die Quelle darf global sein, das Ziel
      ist immer lokal, Punkt 9). Ein fehlender Wert, ein leerer Wert und ein Wert, der nur aus Leerraum
      besteht, sind abgelehnt; die Meldung nennt, welcher der beiden fehlt (beide fehlen: beide in einer
      Zeile, Reihenfolge `user.name`, `user.email`). Abgelehnt ist außerdem ein Wert, der ein Steuerzeichen
      (einschließlich Tabulator, Zeilenumbruch, `CR`) oder ein Byte enthält, das keine gültige UTF-8-Folge
      bildet; die Meldung nennt den Namen der Einstellung, **nie den Wert** (Grund: der Wert kann ein Name
      einer Person sein, und die Meldung geht in Logs).
   4. **Remote** (lokal): `PGR_REMOTE` fehlt (dann legt Punkt 7 es an) oder hat genau die URL aus Punkt 7
      (zeichenweiser Vergleich). Jede andere URL ist abgelehnt; die Meldung nennt beide URLs. Nichts wird
      gesetzt oder entfernt (Randform 3).
   5. **Lesezug im Container** (ein `ssh`-Aufruf, Punkt 18): Antwortet er nicht (Container gestoppt,
      falscher Port, falscher Schlüssel, unbekannter oder geänderter Hostschlüssel, Zeitgrenze), ist der
      Lauf abgelehnt; die Meldung nennt Adresse, Port und die letzte Ausgabe von `ssh` und verweist auf
      `make container-start` (Randformen 4 und 9). Der Aufruf lautet
      `ssh -i $PGR_DEV/ssh/id_ed25519 -o IdentitiesOnly=yes -o BatchMode=yes -o ConnectTimeout=5
      [$PGR_SSH_OPTS] -p <Port> vscode@127.0.0.1 <Lesebefehl>`; **`BatchMode=yes` steht in jedem Aufruf des
      Ziels** (nie eine Rückfrage, nie ein Warten auf eine Eingabe), und das Ziel setzt **keine** Option für
      den Hostschlüssel (Randform 9, Punkt 20 (5)). Der Lesebefehl schreibt nichts (kein `git init`, kein
      `mkdir`, kein `git fetch`) und gibt Zeilen `<schlüssel>=<wert>` aus; er liefert: `git` und `make` im
      Pfad des Containers; zum Bare-Repo `/workspaces/<Projektname>/austausch.git` fehlt / ist ein Bare-Repo
      mit `HEAD` auf `refs/heads/main` / ist ein Bare-Repo mit anderem `HEAD` / ist etwas anderes; zum Klon
      `/workspaces/<Projektname>/pgwire-recorder` fehlt / ist ein Git-Repo mit `origin` gleich dem Pfad des
      Bare-Repos und ohne uncommittete Änderung (`git status --porcelain` leer) / ist ein Git-Repo mit
      uncommitteter Änderung / ist ein Git-Repo mit anderem `origin` / ist etwas anderes; ob der Zweig
      `PGR_ZWEIG` im Klon besteht.
   6. **Werkzeuge im Container**: fehlt `git` oder `make`, ist der Lauf abgelehnt; die Meldung nennt das
      Werkzeug und den Container (Grund: beide liegen im Image, und ein Eingriff in das Image ist nicht
      Gegenstand dieses Ziels).
   7. **Bare-Repo**: `HEAD` zeigt nicht auf `refs/heads/main` oder der Pfad ist kein Bare-Repo —
      abgelehnt, die Meldung nennt den Pfad und den gefundenen `HEAD` (Randform 1; Grund: ein Bare-Repo mit
      `HEAD` auf `master` ließe `git clone` ohne Arbeitskopie enden, Dokument §6).
   8. **Klon**: uncommittete Änderung, anderer `origin`, oder der Pfad besteht und ist kein Git-Repo —
      abgelehnt, die Meldung nennt den Pfad und den Grund und, bei uncommitteter Änderung, dass der
      Aufrufer committet oder stasht (Randform 2; Grund: das Ziel überschreibt keine Arbeit).

   Fälle `p5-*`, je Zeile und je Element einer; Mutationen: die Prüfung entfällt, die Reihenfolge zweier
   Prüfungen ist vertauscht, der Lesezug schreibt, `BatchMode` fehlt, `ConnectTimeout` fehlt, mehr als ein
   `ssh`-Aufruf in `zustand`, die Meldung nennt den Wert der Identität.

6. **Bare-Repo.** Schritt `bare`. Fehlt es, legt das Ziel es an:
   `git init --bare --initial-branch=main /workspaces/<Projektname>/austausch.git` im Container, per `ssh`.
   Ist es nach Punkt 5 vorhanden und in Ordnung, **folgt kein Aufruf** (kein zweites `git init`, keine
   Zeile `Reinitialized`; Randform 1). Ein Fehler endet den Lauf mit **Exit-Code 5** und der Meldung `bare`.
   Fälle `p6-*`: fehlt (Aufruf mit beiden Argumenten, danach `HEAD` auf `main`), vorhanden (kein Aufruf im
   Protokoll), Fehler des Aufrufs (Exit 5, kein Folgeschritt); Mutationen: `--initial-branch` fehlt,
   `git init` läuft auch auf dem vorhandenen Repo, der Fehler wird überlesen.

7. **Remote und Einlieferung.** Schritt `einliefern`, am Host. Fehlt `PGR_REMOTE`, legt das Ziel es an:
   `git remote add <PGR_REMOTE> ssh://vscode@127.0.0.1:<PGR_SSH_PORT>/workspaces/<Projektname>/austausch.git`
   — diese URL ist die zugesagte Form (wörtlich, mit dem Port aus der Eingabe und dem Projektnamen aus
   Punkt 3; Randform 8: je Instanz ein eigener Port und ein eigenes `PGR_REMOTE`, der Pfad trägt den
   Projektnamen). Besteht es mit genau dieser URL, bleibt es unverändert. Dann liefert das Ziel ein:
   `git push <PGR_REMOTE> refs/heads/main:refs/heads/main`, mit `GIT_SSH_COMMAND` gleich
   `ssh -i $PGR_DEV/ssh/id_ed25519 -o IdentitiesOnly=yes -o BatchMode=yes [$PGR_SSH_OPTS]` in der Umgebung
   genau dieses Aufrufs. **`core.sshCommand` des Host-Repos liest und schreibt das Ziel nie** (Grund: die
   Einstellung gilt dem ganzen Repo, nicht einem Remote (Dokument §9), und bei mehreren Instanzen trüge sie
   den Schlüssel der falschen; `GIT_SSH_COMMAND` gilt genau diesem Aufruf). Kein `--force`, kein
   `--force-with-lease`, kein anderer Refspec. Wird der Push abgewiesen (Bare-Repo hat einen `main` mit
   anderer Historie, `! [rejected]`), endet der Lauf mit **Exit-Code 5**; die Meldung nennt `einliefern`,
   den Remote und `git fetch <PGR_REMOTE>` als Weg des Aufrufers; nichts wird zurückgebaut und nichts
   erzwungen (Randform 13). Fälle `p7-*`: Remote fehlt (Aufruf, URL wörtlich), Remote vorhanden (kein
   `remote add` im Protokoll), Port abweichend in der URL, Push-Aufruf mit Refspec und ohne `--force`,
   `GIT_SSH_COMMAND` in der Umgebung des Pushes (Protokoll der Attrappe mit Umgebung) und **nicht** in der
   Konfiguration, Push abgewiesen (Exit 5, Meldung, kein zweiter Versuch), Push mit Fehler ohne
   Abweisung (Exit 5); Mutationen: `--force` ergänzt, Refspec `main` ohne `refs/heads/`, URL ohne Port,
   `core.sshCommand` wird gesetzt, der Fehler wird überlesen, zweiter Versuch nach Abweisung.

8. **Klon.** Schritt `klon`, im Container. Fehlt er, legt das Ziel ihn an:
   `git clone /workspaces/<Projektname>/austausch.git /workspaces/<Projektname>/pgwire-recorder`. Besteht
   er nach Punkt 5 in Ordnung, führt das Ziel allein `git -C <Klon> fetch --prune origin` aus: **kein
   `reset`, kein `checkout`, kein `merge`, kein `pull`, kein `switch` außer dem aus Punkt 11** (Randform 2;
   Grund: der Klon kann Arbeit tragen, und der Stand, aus dem Punkt 11 den Zweig bildet, ist
   `origin/main`). Ein Fehler endet den Lauf mit **Exit-Code 5** und der Meldung `klon`. Fälle `p8-*`:
   fehlt (Aufruf, danach `origin` gleich dem Bare-Pfad und ein Arbeitsbaum), besteht (nur `fetch` im
   Protokoll), Fehler des Aufrufs; Mutationen: `clone` läuft auch auf dem vorhandenen Klon, `fetch` wird
   durch `pull` ersetzt, `--prune` fehlt, der Fehler wird überlesen.

9. **Identität des Klons.** Schritt `identitaet`, im Container. Das Ziel setzt
   `git -C <Klon> config --local user.name <Wert>` und `--local user.email <Wert>` mit den Werten aus Punkt
   5 (3). **Nie `--global`, nie `--system`, nie ein Name oder eine Adresse aus dem Skript oder dem
   Dokument** (Zusage; Fall: `~/.gitconfig` und `/etc/gitconfig` des Containers sind nach dem Lauf
   unverändert, und kein Aufruf im Protokoll trägt `--global` oder `--system`; Mutationen: je eine Ebene
   gesetzt, ein Wert aus dem Skript).

   **Transport.** Die beiden Werte reisen über den **Standardeingang** des `ssh`-Aufrufs, je einer in einer
   Zeile, und werden im Container mit `IFS= read -r` gelesen und als Argument an `git config` gegeben;
   **nie steht ein Wert in der Kommandozeile eines `ssh`-Aufrufs und nie in einer Zeichenkette, die eine
   Shell auswertet** (Grund: `printf %q` am Host setzt die Regeln der Shell des Hosts voraus, und die Shell
   im Container ist die Anmelde-Shell von `vscode`; über den Standardeingang gelesen, ist der Wert ein
   Datum und kein Programm). Zeichen, die unverändert ankommen: Leerzeichen, `"`, `'`, `` ` ``, `$`, `\`,
   `%`, `;`, `&`, `|`, `<`, `>`, `*`, `?`, `#`, `!`, `(`, `)`, `{`, `}`, `[`, `]`, `~`, Nicht-ASCII in
   UTF-8 (`ä`, `日`, Emoji). Abgelehnt sind Steuerzeichen und Nicht-UTF-8, schon in Punkt 5 (3). Fälle
   `p9-*`: je Zeichen der Aufzählung einer (der Wert kommt byteweise gleich an, und das Protokoll der
   Attrappen nennt keinen Aufruf, den der Wert erzeugt hätte — etwa `id` aus `$(id)` oder `` `id` ``),
   beide Werte gesetzt, Reihenfolge der zwei Zeilen; Mutationen: der Wert reist als Argument, `read` ohne
   `-r`, `read` ohne `IFS=`, die Werte vertauscht, der Wert wird in eine Shell-Zeichenkette gesetzt.

10. **Commit-Träger.** Schritt `hooks`: `make hooks-install` im Klon (`cd <Klon> && make hooks-install`
    per `ssh`). Ein Fehler endet den Lauf mit **Exit-Code 5** und der Meldung `hooks`. Der Schritt läuft
    bei jedem Aufruf, auch beim zweiten (das Ziel ist idempotent, `core.hooksPath` gilt je Klon). Fälle
    `p10-*`: Aufruf im Protokoll der Attrappe von `make` mit Ziel und Arbeitsverzeichnis, Fehler (Exit 5,
    kein Folgeschritt); Mutationen: anderes Ziel, anderes Verzeichnis, Fehler überlesen.

11. **Zweig.** Schritt `zweig`, im Container. Besteht `PGR_ZWEIG` im Klon nicht, legt das Ziel ihn an:
    `git -C <Klon> switch -c <PGR_ZWEIG> origin/main`. Besteht er, wechselt das Ziel auf ihn:
    `git -C <Klon> switch <PGR_ZWEIG>` — **ohne** `-C`/`--force`, ohne `reset`, ohne `rebase`, ohne
    `--discard-changes` (Randform 16: Wiederholung nach einem Neustart führt die Arbeit fort). Quelle ist
    immer `origin/main` nach dem `fetch` aus Punkt 8, nie der lokale `main` des Klons (Grund: der lokale
    `main` kann hinter dem Bare-Repo stehen, und das Ziel bewegt ihn nie). Ein Fehler endet den Lauf mit
    **Exit-Code 5** und der Meldung `zweig`. Fälle `p11-*`: Zweig fehlt (Aufruf mit `-c` und `origin/main`,
    danach ist er ausgecheckt), Zweig besteht (Aufruf ohne `-c`, kein `reset` im Protokoll), lokaler `main`
    hinter `origin/main` (der neue Zweig trägt den Stand von `origin/main`), Fehler des Aufrufs; Mutationen:
    Quelle `main` statt `origin/main`, `-c` auch beim vorhandenen Zweig, `--force` ergänzt, der lokale
    `main` wird bewegt.

12. **Träger des Harness, letzter Schritt.** Schritt `traeger`: `make traeger-fetch` im Klon per `ssh`,
    **nach** `zweig`. **Grund für die Stelle am Ende** (Randform 10): Der Schritt braucht Netz im
    Container, und das Egress des Gerüsts ist eingeschränkt (`SPEC-058` Punkt 7). Scheitert er, stehen
    Bare-Repo, Klon, Identität, Träger-Pfad und Zweig schon, und die Wiederholung ist ein zweiter Aufruf
    des Ziels. Ein Fehler endet den Lauf mit **Exit-Code 5** und der Meldung `traeger`; nichts wird
    zurückgebaut, und die Schlussmeldung aus Punkt 13 bleibt aus. Fälle `p12-*`: Aufruf im Protokoll
    (Ziel, Verzeichnis), Fehler (Exit 5, keine Schlussmeldung, die Schritte davor sind gelaufen), Stelle
    (das Protokoll nennt `traeger-fetch` **nach** dem `switch` aus Punkt 11); Mutationen: der Schritt läuft
    vor `zweig`, der Fehler wird überlesen, nach dem Fehler erscheint die Schlussmeldung.

13. **Erfolg, Wiederholbarkeit, Zugang.** Erfolg ist eine Zeile auf stdout
    `container-einrichten: eingerichtet <Klon-Pfad> <Zweig>` und **Exit-Code 0**, beim ersten wie beim
    zweiten Aufruf. Jeder Schritt ist wiederholbar (Punkte 5 bis 12); ein zweiter Aufruf mit denselben
    Eingaben ändert nichts außer `fetch`, Identität, `hooks-install`, `switch` und `traeger-fetch` und
    endet mit derselben Zeile (Fall `p13-zweiter-aufruf`: zwei Läufe nacheinander, der zweite ohne
    `git init`, ohne `clone`, ohne `remote add`, mit Exit 0; Mutation: der zweite Lauf legt an).
    **Akzeptierte Negative, je mit Grund:** (1) Bricht der Lauf während `git init --bare` ab, kann ein
    Verzeichnis ohne `HEAD` auf `main` bleiben; der nächste Aufruf lehnt es nach Punkt 5 (7) ab und nennt
    den Pfad — ein Rückbau wäre ein Mechanismus für ein Zeitfenster von Millisekunden in einem Werkzeug,
    das nicht das Ziel ist. (2) Gleichzeitige Aufrufe mit demselben `PGR_DEV` sind nicht zugesagt (Grund
    wie `SPEC-058` Punkt 17 (6): ein Nutzer, ein Terminal; eine Sperre bräuchte ein Aufräumen nach
    `SIGKILL`). (3) Dass der eingelieferte `main` der Stand ist, den der Aufrufer meint, prüft das Ziel
    nicht; es liefert `refs/heads/main` ein, wie er ist (uncommittete und ignorierte Dateien kommen nicht
    mit, Dokument §9).

    **Zugang von Claude Code: nichts davon** (Randform 25, Bindung an
    `slice-harness-devcontainer-einrichten-zugang`). Das Ziel liest, löst, prüft und übergibt keinen Zugang
    und nennt keine der Variablen `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GIT_TOKEN` in einem
    Aufruf; es **entfernt sie zu Beginn aus seiner eigenen Umgebung**, sodass kein Kindprozess (`ssh`,
    `git`, `make`) sie erbt (wie `SPEC-058` Punkt 8; Fall `p13-umgebung` mit einem Sentinel in allen drei
    Variablen: keine Attrappe sieht eine, keine Ausgabe und keine Datei trägt den Wert; Mutation: je
    Variable eine, die das Entfernen entfällt). **Die Stelle für die Erweiterung** ist ein Schritt
    `zugang` **vor** `eingaben`: Der Ablauf in Punkt 4 ist eine numerierte Folge, in die er passt, ohne
    dass eine bestehende Zusage sich ändert; wer ihn einfügt, ergänzt Punkt 4 und Punkt 14 um ihn.
    **Empfehlung des Architect an den Nehmer:** Braucht die Auflösung einen Wert aus Punkt 2, steht sie
    hinter `eingaben` statt davor; das ist dann eine Änderung dieses Punktes mit eigenen Fällen.

14. **Reihenfolge der Gründe und Ausgang.** Treffen mehrere Gründe zu, gilt diese Reihenfolge, der erste
    beendet den Lauf (innerhalb von `eingaben` und `werkzeuge` werden alle Verstöße des Schrittes
    gemeldet, Punkte 2 und 3):

    1. `eingaben` (Exit 2) · 2. `werkzeuge` am Host (3) · 3. `zustand` (4) in der Reihenfolge der Liste in
    Punkt 5: Gerüst, Host-Repo, Identität, Remote, Lesezug, Werkzeuge im Container, Bare-Repo, Klon ·
    4. Fehler eines Schrittes `bare`, `einliefern`, `klon`, `identitaet`, `hooks`, `zweig`, `traeger` (5) ·
    5. Signal (130, 143).

    **Die Identität steht vor dem Container** (Randform 5 gegen Randform 4): Sie ist lokal und ohne
    Container prüfbar, und der Aufrufer behebt sie ohne laufenden Container.

    **Eine Ausnahme von „Eingabe vor Werkzeug": die Prüfung, die das Werkzeug selbst ist** (Entscheidung
    des Architect, 2026-10-21). Ist `PGR_ZWEIG` allein nach `git check-ref-format --branch` ungültig und
    fehlt zugleich `git` am Host, bleibt diese eine Prüfung aus, und der Lauf endet im Schritt `werkzeuge`
    mit **Exit 3** (`git fehlt im Pfad`), nicht mit Exit 2. **Grund:** Die Reihenfolge kann keine Prüfung
    vorziehen, die das fehlende Werkzeug braucht; die Gegenrichtung — `eingaben` lehnt ohne `git` jeden
    Zweignamen ab — meldete dem Aufrufer den Zweignamen, während sein Mangel das fehlende `git` ist, und
    machte den Fall `p3-fehlt-git` mit gültigem Namen unerreichbar. **Die Ordnung selbst bleibt:** Jede
    Prüfung aus Punkt 2, die ohne `git` auskommt — `PGR_DEV`, `PGR_SSH_PORT`, `PGR_REMOTE`, `PGR_SSH_OPTS`
    und für `PGR_ZWEIG` die Formen nicht gesetzt, leer, `main`, führendes `-`, `refs/…` —, meldet auch ohne
    `git` mit Exit 2 (Fall `p14-eingabe-vor-werkzeug`).

    **Signal und Fehler eines
    Schrittes sind Ereignisse, keine Zustände:** Das erste beendet den Lauf (Punkt 15). Exit-Code 0 heißt
    „eingerichtet"; 1 und alle anderen sind nicht belegt. Fälle `p14-*`: je Nachbarpaar der Reihenfolge
    ein Fall, in dem beide Gründe zutreffen und der frühere gemeldet wird (Eingabe vor Werkzeug; Werkzeug
    vor Gerüst; Gerüst vor Host-Repo; Host-Repo vor Identität; Identität vor Remote; Remote vor Lesezug;
    Lesezug vor Werkzeugen im Container; Werkzeuge vor Bare-Repo; Bare-Repo vor Klon; Klon vor Fehler eines
    Schrittes), dazu der Exit-Code jedes Grundes; Mutationen: je Paar die beiden Prüfungen vertauscht.

15. **Abbruch.** `SIGINT` und `SIGTERM` beenden den Lauf: Das Skript beendet seinen laufenden Kindprozess
    (`ssh`, `git`), meldet auf stderr `container-einrichten: <Schritt>: abgebrochen` mit dem Schritt, in
    dem es war, und endet mit **130** beziehungsweise **143**. Es baut nichts zurück (Randform 11); jeder
    Schritt ist wiederholbar (Punkt 13). `SIGKILL` und ein Ausfall des Hosts fangen nichts auf. Fälle
    `p15-sigint`, `p15-sigterm` an einer Attrappe von `ssh`, die wartet (Exit-Code, Meldung mit dem
    Schritt, der Kindprozess lebt nicht mehr), und `p15-wiederholung` (der nächste Aufruf setzt fort);
    Mutationen: die Falle beendet das Kind nicht, die Falle entfällt für `INT` oder für `TERM` (äquivalent,
    wenn bash allein 130 beziehungsweise 143 ergibt: die Gegenprobe führt es dann als Mutant ohne roten
    Fall, wie `SPEC-057` Punkt 7; gemessen vom Implementer, mit Invariante, Ort und Sonde nach
    `AGENTS.md` §3.10).

16. **Aufrufarten des Treibers.** Das Ergebnis ist der Exit-Code des Skripts und die Zeile `eingerichtet`;
    kein Stand und kein Hook glaubt es. Das Rezept von `container-einrichten` ist eine Zeile **ohne
    `$(MAKE)` und ohne `+`**, das Ziel `.PHONY` ohne Voraussetzungen. Unter `-n` zeigt `make` das Rezept
    und führt es nicht aus: kein `ssh`, kein `git`, kein Remote, keine Konfiguration, keine Zeile
    `eingerichtet`. Unter `-q` und `-t` ebenso nichts, und kein Zeitstempel, weil das Ziel keine Datei ist.
    Unter `-i` überliest `make` den Exit-Code des Rezepts: Der Exit-Code von `make` belegt dann nichts, und
    die Zusage ist die Meldung — `container-einrichten: eingerichtet …` steht nur nach einem Erfolg, bei
    jedem Fehler stattdessen die Zeile des Schrittes auf stderr. `MAKEFLAGS` aus der Umgebung gilt wie das
    Flag, auch in Kombination; bei `-n`, `-q` oder `-t` zusammen mit `-i` gilt `-n`, `-q`, `-t` (nichts
    läuft). `-k` und `-j` ändern an diesem einen Ziel nichts: gleiches Ergebnis wie ohne. Je Flag `-n`,
    `-q`, `-t`, `-i`, `-k`, `-j`, je `MAKEFLAGS` aus der Umgebung mit demselben Flag und für `-n` zusammen
    mit `-i` ein Fall `p16-*` an Attrappen (Protokoll leer bei `-n`, `-q`, `-t`; bei `-i` ein Fehlerfall
    ohne Zeile `eingerichtet`; bei `-k` und `-j` ein Erfolgsfall); Mutation: ein Rezept mit `$(MAKE)` oder
    `+`.

17. **Gegenprobe: Attrappen statt Container, und warum.** `container-einrichten-gegenprobe` ist `.PHONY`,
    steht in `GEGENPROBEN_CONTAINER` (sie ruft weder `u-boot` noch `devcontainer`,
    [`SPEC-060`](SPEC-060-gegenproben.md) Punkte 2 und 3), hängt nicht an `GATE_CHECKS` und steht nicht in
    der Ausnahmeliste der schnellen Kette. Sie läuft in Temp-Verzeichnissen (`mktemp -d`) **ohne Container,
    ohne Docker, ohne Netz und ohne echtes `ssh`**:

    - **Attrappe von `ssh`** in einem vorangestellten Pfad: protokolliert Argumente und Standardeingang und
      führt den übergebenen Befehl in einer lokalen Shell in der Welt des Falls aus (die Welt enthält
      `workspaces/<Projektname>/` als Verzeichnis). So laufen `git init --bare`, `git clone`,
      `git config --local`, `git switch` und `git status` als **echtes `git`** auf lokalen Pfaden, und die
      Zusagen der Punkte 5 bis 11 sind an echtem `git`-Verhalten geprüft (`Reinitialized`, Exit 128 bei
      nicht leerem Ziel, `! [rejected]`), nicht an einer Nachbildung.
    - **Attrappe von `make`** für die Schritte `hooks` und `traeger`: protokolliert Ziel, Verzeichnis und
      Umgebung und liefert den Exit-Code, den der Fall vorgibt.
    - **Welt je Fall:** ein Host-Repo mit **einem** Commit (`git init`, ein Commit, Identität lokal), das
      Verzeichnis `$PGR_DEV` mit `u-boot.yaml` und `ssh/id_ed25519` (eine Datei, kein Schlüsselpaar: das
      Ziel liest sie nur), und das Ziel-Verzeichnis `workspaces/<Projektname>/`.
    - **Aufräumen:** Der Ordner eines Falls entfällt, sobald der Fall endet, der einer Mutation, sobald
      ihre Fälle bewertet sind — rot wie grün ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16 (5)); bei
      `SIGINT`, `SIGTERM` und Ende entfällt der ganze Arbeitsordner.
    - **Zeitbudget:** Die Gegenprobe ist **ab ihrem Anlegen aufgenommen** im Sinn von
      [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16: Sie trägt die Fälle und Mutationen von dessen Punkt
      16 (9) (`GP_BUDGET`, Zeile mit Zeit und Last, Überschreitung rot, Reihenfolge der Gründe, Ordner
      entfernt), und der Slice, der sie anlegt, trägt sie in der Tabelle dort als `aufgenommen` ein und
      nennt die gemessene Zeit. Die Obergrenze ist 1200 s
      ([ADR-0059](../../docs/plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), vorher 600 s nach
      [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)).
    - **Start:** Die Gegenprobe setzt `SIGINT` beim Start als nicht ignoriert voraus. Ein Start mit ignoriertem `SIGINT` (Hintergrundstart mit `&` oder
      `nohup`, auch `setsid -f` allein) ist für ihre Fälle mit `SIGINT` nicht zugesagt. Aufruf im Vordergrund oder mit
      `setsid -f env --default-signal=INT,QUIT`; Rezept, Messung und Grund stehen in
      [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 17 (akzeptiertes Negativ, ohne Fall).
    - **Grün ohne vorherigen Lauf von `make gates`** ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 3, Regel
      an die Lieferung): Das belegt der Slice, der die Gegenprobe anlegt, durch **ihren eigenen** Lauf in
      einem frischen Gegenproben-Container — nicht durch eine Fahrt des Sammelziels (Entscheidung des
      Architect, 2026-10-21; Grund: Die Zusage gilt dieser Gegenprobe, und ihr eigener Lauf belegt sie
      vollständig; eine Fahrt von `make gegenproben` fährt heute die Gegenproben mit, die noch über der
      Obergrenze liegen, und ist deshalb nach [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)
      *Übergang* vorerst keine Handlung eines Slice). Die Pflicht, `make gegenproben` und
      `make gegenproben-host` am Ende einer Welle und vor einem Release zu fahren (`AGENTS.md` §6,
      `SPEC-060` Punkt 14), bleibt davon unberührt; dieser Vertrag ändert sie nicht.

    **Grund, warum kein echter Container Teil der Gegenprobe ist** (Entscheidung des Architect, 2026-10-21,
    innerhalb von [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)): Gemessen im
    Dev-Container (20 Kerne, Last 4,8) kostet ein Fall mit Attrappen 25 bis 60 ms (Welt 14 ms, Gutfall aus
    `init --bare`, `push`, `clone`, zwei `config` und `switch` 43 ms, Abbruchfall 8 ms, Aufräumen 3 ms);
    ein Fall gegen einen echten Container kostet 3 bis 25 s (vier bis sechs `ssh`-Runden zu je 0,1 bis
    0,3 s, `make` im Container, `traeger-fetch` mit Netz), und der Start eines Scratch-Gerüsts kostet
    einmalig 65 s bei warmem Cache (`slice-harness-devcontainer-u-boot-0-10`). Eine Mutation fährt die
    Fälle ihrer Gruppe erneut; bei rund 65 Fällen und rund 120 Mutationen ergibt das mit Attrappen 200 bis
    400 s bei `GP_PARALLEL` 4 und mit einem echten Container ein Vielfaches der Obergrenze — die
    vergleichbare `container-start-gegenprobe` braucht mit ihrer Gruppe Docker **19 min 38 s**
    (`slice-harness-devcontainer-einrichten-start` §7). Dazu käme, dass eine Gegenprobe mit Container nach
    `SPEC-060` Punkt 3 in die Liste für den Host gehörte und damit weder in `make gegenproben` noch im
    Dev-Container liefe, in dem der Implementer arbeitet. Die Zusagen dieses Vertrags hängen an `git`,
    `ssh`-Argumenten, Reihenfolge und Meldungen — alle an Attrappen prüfbar; was ein Container zusätzlich
    belegt, belegt der Nachweis aus Punkt 18.

    **Akzeptiertes Negativ, mit Grund:** Die Gegenprobe belegt nicht, dass ein echtes `ssh` die Argumente
    so annimmt, dass der Befehl im Container läuft, und nicht, dass `make hooks-install` und
    `make traeger-fetch` im Klon gelingen. Grund: Beides braucht einen Container und steht damit über der
    Obergrenze; der Nachweis aus Punkt 18 belegt es einmal je Lieferung, und der Slice trägt den Beleg in
    seiner Closure-Notiz.

18. **Nachweis gegen einen laufenden Container — Handlung, kein Sensor.** Was die Gegenprobe nicht belegt
    (Punkt 17), belegt dieser Nachweis: dass das Ziel gegen einen echten Container durchläuft. Er ist
    **kein Gate, keine Gegenprobe und steht in keiner der beiden Listen** von
    [`SPEC-060`](SPEC-060-gegenproben.md); er hat keinen Fall, keine Mutation und keine
    Abdeckungs-Deklaration (gleiche Lesart wie [`SPEC-067`](SPEC-067-release.md) Punkt 20 (7)): Es urteilt
    niemand außer dem Menschen, der die Ausgabe liest, und kein Hook und kein Stand glaubt sie. Er hat
    **keinen eigenen Treiber** (Grund: jeder Schritt ist schon ein `make`-Ziel oder ein Kommando; ein
    Skript darüber wäre ein weiteres Artefakt mit eigenen Zusagen, ohne einen Schritt hinzuzufügen).

    **Die Schritte, in dieser Reihenfolge, am Host** (`<n>` ist ein freier Port, nie 2222 und nie der Port
    eines laufenden Containers; der Name trägt das Präfix `pgr-ce-`):

    1. `PGR_DEV=$(mktemp -d)/pgr-ce-nachweis PGR_NAME=pgr-ce-nachweis PGR_SSH_PORT=<n> PGR_CPUS=2
       PGR_MEM=2g make container-start` — Gerüst und Container (`SPEC-058`).
    2. `PGR_DEV=… PGR_SSH_PORT=<n> PGR_REMOTE=pgr-ce-nachweis PGR_ZWEIG=<Zweig des Slice> make
       container-einrichten` — der **erste** Aufruf (Bare-Repo und Klon fehlen).
    3. derselbe Aufruf erneut — der **zweite** Aufruf (Bare-Repo, Klon, Remote und Zweig bestehen,
       Randformen 1 bis 3 und 16).
    4. Ablesen im Container, über `ssh` mit dem Schlüssel des Gerüsts, je ein Kommando mit seiner Ausgabe:
       `git -C <Bare> symbolic-ref HEAD`; `git -C <Klon> remote get-url origin`;
       `git -C <Klon> rev-parse --abbrev-ref HEAD`; `git -C <Klon> config --local --get user.name` und
       `user.email` (gleich den Werten des Host-Repos); `git -C <Klon> config --get core.hooksPath`;
       `ls <Klon>/.harness/` oder die Meldung von `make traeger-fetch`.
    5. **Aufräumen und seine Ausgabe, nur per Name und nur Objekte des Präfixes, nie `prune`** (Lesart von
       `SPEC-058` Punkt 15 (d)): `docker rm -fv pgr-ce-nachweis` entfernt den Container und seine
       **anonymen** Volumes; die zwei **benannten** Volumes des Gerüsts (`pgr-ce-nachweis-workspace-…`,
       `pgr-ce-nachweis-docker-…`) und das Image-Paar (`vsc-pgr-ce-nachweis-…`, `…-uid`) überleben ihn und
       werden danach mit `docker volume rm` und `docker rmi` per Name entfernt. Zuletzt
       `docker ps -a --filter name=pgr-ce-`, `docker volume ls` und `docker images` (kein Container, kein
       Volume und kein Image des Präfixes bleibt), `git remote remove pgr-ce-nachweis` im Host-Repo und
       `rm -rf` des Temp-Verzeichnisses. Der Beleg zeigt **beide** Zustände: die Liste nach
       `docker rm -fv` (die zwei Volumes und die zwei Images stehen noch) und die Liste nach dem Entfernen
       per Name (keines bleibt).

    **Beleg** (Ort und Form wie [`SPEC-067`](SPEC-067-release.md) Punkt 20 (6), *Ort des Belegs*, Entscheidung des Nutzers
    vom 2026-10-22): eine **eigene Datei** `docs/reviews/<YYYY-MM-DD>-nachweis-<Kennung des Slice>.md`. Ihr Kopf nennt Datum,
    `uname -s`/`uname -m`, den Stand (Commit) und die Versionen von `docker version --format '{{.Server.Version}}'`, `u-boot --version`,
    `devcontainer --version` und `git --version` des Hosts; danach stehen die Schritte 1 bis 5 mit ihren Ausgaben **ungekürzt**, darunter
    beide Exit-Codes aus 2 und 3 und die Zeile `container-einrichten: eingerichtet …`. **§7 des Slice-Plans nennt** den Pfad der Datei,
    ihre Summe nach `sha256sum`, die Zeile `container-einrichten: eingerichtet …` mit den beiden Exit-Codes und den Kopf (Datum, Stand,
    Plattform); die Ausgabe selbst steht nicht in §7. Gekürzt wird allein ein Geheimnis oder ein kurzlebiger Token, mit der Marke
    `[gekürzt: <Grund>]` (Lesart von `SPEC-067` Punkt 20). Eine Angabe ohne diese Datei ist eine Angabe der Koordination, kein Beleg; die
    Zusagen aus Punkt 17 (*akzeptiertes Negativ*) gelten dann als nicht geprüft. [`ADR-0057`](../../docs/plan/adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md)
    nennt in der Tabelle seiner Fitness Function noch „Beleg in §7 des Slice-Plans“; maßgeblich ist dieser Punkt (Rang 2 vor Rang 4), und die ADR wird dafür
    nicht ersetzt (akzeptiertes Negativ; Grund: Die Zelle wiederholt eine Regel der Spezifikation, `AGENTS.md` §3.8, und eine Ersetzungs-ADR berichtigte
    nur den Text).

    **Akzeptiertes Negativ, mit Grund:** Bricht der Nachweis zwischen Schritt 1 und Schritt 5 ab, bleibt ein
    Container des Präfixes `pgr-ce-` samt seinen Volumes **und seinen Images** stehen. Grund: Ohne Treiber fängt ihn keine Falle;
    das Präfix macht ihn auffindbar, und Schritt 5 steht mit seiner Ausgabe im Beleg, sodass ein fehlendes
    Aufräumen im Beleg sichtbar ist. Der Nachweis läuft nie gegen den Container, in dem gearbeitet wird,
    und nie gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon (Regel an den Aufrufer,
    `SPEC-052` Punkt 8, `SPEC-058` Punkt 15).

19. **Was die Plattform entscheidet.** Geprüft ist Linux (Host und Container) mit bash ≥ 4.4 am Host
    (Kopf dieser Kennung, *Plattform*), und die Gegenprobe läuft unter Linux im Gegenproben-Container. Für
    bash 4.0 bis 4.3 gilt das akzeptierte Negativ im Kopf (*Warum 4.4*), ohne Selbsttest des Ziels. Für
    macOS und Windows ist **keine** Zusage dieses Vertrags
    geprüft; sie sind die **Grenze** (Kopf dieser Kennung). Je übriger Plattform die Entscheidung: Das Ziel
    lehnt sie nicht ab und wird dort nicht gemessen (Grund: kein solcher Host steht bereit, und das
    Gerüst aus `SPEC-058` ist dort ebenso nicht gemessen); ein Satz im Dokument
    `docs/developer/devcontainer.md` über eine dieser Zusagen nennt allein Linux (`AGENTS.md` §3.11). Der Abbruch (Punkt 15) ist allein unter Linux geprüft; der Start der Gegenprobe mit
    ignoriertem `SIGINT` ist auf keiner Plattform zugesagt (Punkt 17, [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 17).

20. **Grenzen.** (1) Versionen von `ssh` und `git` am Host und von `git` und `make` im Container: nicht
    geprüft, nur das Vorhandensein (Punkte 3 und 5); ein `git` im Container ohne `--initial-branch` endet
    als Fehler des Schrittes `bare` (Exit 5). (2) Was im Container außer Bare-Repo, Klon, Identität,
    `core.hooksPath` und Zweig liegt, ist nicht Gegenstand; der Inhalt des Klons jenseits dieser Größen wird
    nicht verglichen. (3) Die Netzlage des Containers (Egress-Liste, DNS) ist Eigenschaft des Gerüsts
    (`SPEC-058` Punkt 7) und zeigt sich in diesem Ziel allein als Fehler des Schrittes `traeger`; die
    Version von `u-boot`, die das Gerüst erzeugte, ist für dieses Ziel ohne Belang, weil es `u-boot` nie
    aufruft. (4) Ob `dockerd` im Container bereit ist (`SPEC-058` Punkt 17 (4)). (5) Der Hostschlüssel: Das
    Ziel setzt keine Einstellung und ändert keine Datei; ein unbekannter oder nach einem Neubau geänderter
    Hostschlüssel endet mit Exit 4 und der Ausgabe von `ssh` (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`).
    Wer eine eigene `known_hosts`-Datei braucht, gibt sie in `PGR_SSH_OPTS` (Punkt 2). (6) Gleichzeitige
    Aufrufe mit demselben `PGR_DEV` (Punkt 13). (7) Eine Option in `PGR_SSH_OPTS`, die Leerraum enthält, und eine, die ein
    Shell-Metazeichen enthält: In den `ssh`-Aufrufen des Ziels reist der Wert als einzelne Argumente, die
    keine Shell sieht (Punkt 2, Fall `p2-opts-keine-shell`); im Push reist er nach Punkt 7 in
    `GIT_SSH_COMMAND`, und `git` gibt diese Zeichenkette an eine Shell — dort wird ein Metazeichen
    ausgewertet. Was dabei geschieht, ist **nicht zugesagt und nicht gemessen** (Entscheidung des
    Architect, 2026-10-21). **Grund:** Die beiden Formen — Argumentliste hier, Zeichenkette dort — gibt die
    Schnittstelle von `ssh` und `git` vor, nicht dieses Ziel; der Wert stammt aus der Umgebung des
    Aufrufers am eigenen Host, und eine Ablehnung von Metazeichen schlösse Formen aus, die `ssh` braucht
    (etwa die eigene `known_hosts`-Datei aus (5)). `PGR_DEV` reist in derselben Zeichenkette, trägt aber eine
    geprüfte Form (Punkt 2) und ist deshalb keine Grenze.
    (8) Eine Zeitgrenze für den Lauf insgesamt: Es gibt keine außer `ConnectTimeout=5` je `ssh`-Aufruf; ein
    `git push` oder ein `make traeger-fetch`, das hängt, hängt, und `Strg+C` bricht ab (Punkt 15).
    (9) Ob das Dokument `docs/developer/devcontainer.md` die Schritte richtig nennt, prüft kein Gate.
    (10) Dass der Klon nach dem Lauf einen grünen `make gates` fährt: gemessen von
    `slice-harness-devcontainer-u-boot-0-10`, nicht von diesem Ziel. (11) Der Exit-Code, mit dem
    `git push` eine gescheiterte Verbindung meldet (gemessen 128): Der Vertrag bindet ihn nicht. Zugesagt
    und geprüft ist, was das **Ziel** tut — **Exit 5**, die Meldung
    `einliefern: git push <Remote> endete mit Exit <Zahl>` in dieser Form, kein zweiter Versuch, nichts
    erzwungen (Punkt 7, Fall `p7-fehler`). **Akzeptiertes Negativ, mit Grund** (Entscheidung des Architect,
    2026-10-21): Die Zahl stammt von `git`; sie zu binden sagte das Verhalten eines fremden Werkzeugs über
    seine Versionen hinweg zu, obwohl (1) dessen Version nicht prüft. Die Gegenprobe prüft deshalb die Form
    der Zeile, nicht die Zahl.
