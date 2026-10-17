# SPEC-057 — Aufräumen der Integrationstests (`test-integration`, `integration-volumes-gegenprobe`)

Der Runner `tools/test/run-integration-tests.sh` (Ziel `make test-integration`, auch von
`make test-postgres-versionen` gefahren) startet Container und räumt sie nach dem Lauf auf. Diese
Kennung legt fest, welche Docker-Volumes er dabei mitentfernt und prüft, und
`make integration-volumes-gegenprobe` prüft das. Sie sagt nichts über die Tests selbst, nichts über
Netz und Container-Ende (das tut `SPEC-053` Punkt 7) und nichts über Volumes, die dieser Lauf nicht
angelegt hat.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Die Gegenprobe
fährt gegen reales Docker (über `DOCKER`), legt Container und Volumes unter dem Namenspräfix
`pgr-ivg-` an und entfernt sie per Name; sie ruft weder `prune` auf noch berührt sie ein Objekt
mit anderem Namen. Die unveränderte Kopie ist grün, und je Fall bricht eine Mutation genau ihn.

1. **Gegenstand.** Die Container des Laufs sind `$pg` und `$tests` (Namen `<Kennung>-pg` und
   `<Kennung>-tests`, Kennung `PGR_LAUF` oder `pgr-it-<PID>`). Gezählt werden die Mounts vom Typ
   `volume`, die `docker inspect` an diesen beiden Containern nennt, außer dem benannten Volume
   `$daten` des Laufs (Name `<Kennung>-daten`), das der Runner unter seinem Namen entfernt. Ob ein
   Volume „anonym“ heißt, entscheidet nicht das Skript; maßgeblich ist die Herkunft aus dem
   Container. Bind-Mounts zählen nicht. Ein künftiges benanntes Volume des Runners ist wie `$daten`
   auszunehmen, sonst gilt es als Rest (Punkt 4) und macht den Lauf rot.
2. **Entfernen.** Der Runner entfernt `$pg` und `$tests` mit `docker rm -fv`. Mit `rm -f` allein
   bliebe das Volume zurück, das das Image mit `VOLUME` deklariert (Fall: Mutant ohne `-v`).
3. **Lesen vor dem Entfernen.** Der Runner liest die Namen unmittelbar vor dem Entfernen, für jeden
   Zustand des Containers, auch `created`. Existiert ein Container nicht (Image nicht ladbar,
   Start gescheitert vor dem Anlegen, `$tests` nach `--rm`), gibt es keine Namen: die Liste ist leer,
   ohne Fehler und ohne Meldung.
4. **Prüfung nach dem Entfernen und Rest.** Nach dem Entfernen prüft der Runner, dass keines der
   gelesenen Volumes mehr existiert. Ein Rest wird mit Namen gemeldet (kein Inhalt des Volumes: ein
   PostgreSQL-Datenverzeichnis kann Zugangsdaten tragen), per Namen entfernt (nur diese) und macht
   den Lauf rot: Exit-Code 1, auch bei grünen Tests. Endet der Lauf ohnehin mit einem Code ungleich
   0, behält dieser Vorrang (zum Beispiel 143 nach `SIGTERM`). Gelesen werden nur Namen aus den
   eigenen Containern; die Prüfung entfernt nie ein Volume, dessen Name sie nicht so gelesen hat,
   und zählt nie die Volumes des Hosts (parallele Läufe, auch anderer Projekte, sind damit nicht
   betroffen; gleichzeitige Läufe desselben Runners bleiben nicht zugesagt, `SPEC-052` Punkt 8).
5. **Ende des Laufs.** Der Runner ruft Entfernen, Lesen und Prüfung aus seiner Falle `EXIT`; sie
   läuft bei jedem Ende, das der Runner selbst beendet: grüner und roter Lauf, Fehler des Aufbaus
   (auch Start des Servers gescheitert), `SIGINT` und `SIGTERM`. Die Gegenprobe prüft die
   Funktionen (Punkte 2 bis 4) und das Ende durch `SIGINT` und `SIGTERM` (je ein Fall). Ein
   roter Lauf und ein Fehler des Aufbaus gehen durch dieselbe Falle; einen eigenen Fall dafür
   gibt es nicht (akzeptiertes Negativ; Grund: er brauchte das Testimage und einen Lauf der
   Tests). `SIGKILL` und ein Ausfall des Hosts räumen nichts auf und werden nicht aufgefangen (ein
   späterer Lauf findet die Container eines alten nicht, weil ihr Name die PID trägt). Ist der
   Docker-Dienst beim Aufräumen nicht erreichbar, ist nichts zugesagt (Grund: der Lauf selbst
   kann ohne ihn nicht grün gewesen sein).
6. **Die Funktionen liegen in einer Datei, `tools/test/integration-aufraeumen.sh`.** Runner und
   Gegenprobe laden sie; die Gegenprobe prüft damit den Code des Runners, nicht eine Kopie. Die
   Mutanten ersetzen genau eine Stelle in einer Kopie der Datei oder des Runners; trifft die
   Ersetzung nichts oder mehr als eine Stelle, bricht die Gegenprobe ab.
7. **Fälle der Gegenprobe.** Je Fall die Mutation, die nur ihn bricht:

   | Fall | Zusage | Mutation |
   |---|---|---|
   | `p2-grundfall` | nach dem Aufräumen existiert kein gelesenes Volume (Container aus `docker create` des PostgreSQL-Images) | — |
   | `p2-ohne-v` | ohne `-v` bleibt ein Rest, die Meldung nennt den Namen | `rm -f` statt `rm -fv` |
   | `p3-angelegt` | die Namen werden im Zustand `created` gelesen | Namen erst nach dem Entfernen lesen |
   | `p3-kein-container` | kein Container ergibt eine leere Liste ohne Fehler und Meldung | leere Liste wird ein Fehler |
   | `p4-fremd` | fremde Volumes (ein benanntes, ein anonymes) bleiben | die Prüfung entfernt jedes Volume ohne Container |
   | `p4-rest-rot` | ein Rest macht den Lauf rot (Exit 1), auch bei grünen Tests | der Exit-Code der Prüfung wird verworfen |
   | `p4-vorrang` | ein Lauf-Code ungleich 0 behält Vorrang | die Prüfung setzt immer 1 |
   | `p5-sigterm` | Aufräumen und Prüfung laufen nach `SIGTERM` | die Falle für `TERM` entfällt |
   | `p5-sigint` | Aufräumen und Prüfung laufen nach `SIGINT` | eine Falle für `INT`, die den Runner ohne Aufräumen beendet |

   Die Fälle `p5-*` starten den Runner (Kopie) mit `PGR_BAUEN=0` und `PGR_LAUF` unter dem
   Präfix, warten, bis der Container `$pg` existiert, und senden das Signal; das Testimage wird
   nicht gebraucht. `SIGINT` setzt voraus, dass der Runner nicht als Hintergrundjob einer
   nicht interaktiven Shell startet (sonst wird es ignoriert); das Vorbild ist
   `tools/test/postgres-versionen-gegenprobe.sh`.
8. **Docker-Version.** Gemessen ist Docker 29.9.0: `docker inspect` nennt die Mounts schon im
   Zustand `created`. Eine Version, die das nicht tut, macht `p2-ohne-v` grün statt rot und
   damit die Gegenprobe rot; die Blindheit wird sichtbar. Weitere Versionen sind nicht
   zugesagt; `DOCKER` bleibt der einzige Parameter.
9. **Aufrufarten des Treibers.** Das Ziel schreibt keinen Stand; geglaubt wird allein sein
   Exit-Code in der Kette von `make gates`. Das Rezept ist eine Zeile ohne `$(MAKE)` und ohne
   `+`. Ob `make -n`, `-q`, `-t` oder `-i` das Gate leerlaufen lassen oder überlesen, ist
   nicht Gegenstand dieser Kennung (akzeptiertes Negativ; Grund: kein Hook und kein Stand
   glaubt dieses Gate allein, den Nachweis des vollen Laufs schreibt `make gates` nur ohne diese
   Flags, `SPEC-052` Punkte 6 und 21).
10. **Ausnahmeliste der schnellen Kette.** `integration-volumes-gegenprobe` hängt an
    `GATE_CHECKS` und steht in der Ausnahmeliste (`SPEC-052` Punkte 2 und 4): Es startet Container
    und läuft in der schnellen Kette nicht. Es liest nichts unter `docs/plan/planning/`,
    `docs/reviews/`, `docs/user/`, `docs/developer/` und nicht die `README.md` (`SPEC-052` Punkt 13): Es legt Kopien
    von Dateien unter `tools/test/` in einem Temp-Verzeichnis an und spricht Docker an.
11. **Grenzen.** (1) Ein `docker run` außerhalb des Runners ohne `--rm` und ohne `-v` legt das
    Leck wieder an; kein Gate sieht es (die übrigen `docker run` des Repos tragen `--rm`, belegt
    im Slice). (2) Volumes, die dieser Lauf nicht angelegt hat, entfernt kein Werkzeug des Repos,
    auch nicht die vom Leck vor diesem Vertrag; ein Host wird vom Nutzer bereinigt. (3) Eine Diagnose
    der vollen Platte im Runner gehört nicht zu diesem Vertrag.
