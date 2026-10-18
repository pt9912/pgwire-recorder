# SPEC-062 — Übernahme des Nachweises des vollen Laufs aus einer anderen Umgebung (`nachweis-uebernehmen`)

`make nachweis-uebernehmen NACHWEIS_QUELLE=<verzeichnis>` legt den Nachweis des vollen Laufs
([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 11), den eine andere Umgebung geschrieben hat (ein Worktree, ein Klon im
Dev-Container), im Arbeitsbaum ab, in dem das Ziel läuft, **genau dann**, wenn (a) die Dateiliste dieses Baums, frisch mit
`dateiliste_erzeugen` berechnet, byte-gleich mit der Dateiliste der Quelle ist, (b) die Liste der Quelle unbeschädigt ist
(`dateiliste_pruefen`) und (c) `gates-passed.diffsha` der Quelle byte-gleich `<SHA-256 der Liste>\n` ist. Sonst lehnt es ab und
ändert nichts. Das Ziel ist kein Gate: Es steht nicht in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette, und
`make gates` fährt es nicht. Der Hook `.claude/hooks/stop-require-gates.sh` bleibt unverändert; er liest die Liste und gibt bei
byte-gleicher Liste frei ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 16). Neu ist ein zweiter Schreiber des Nachweises, kein neuer
Leser. Warum diese Lesart des Nachweises gilt, steht in
[ADR-0049](../../docs/plan/adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md); diese Kennung wiederholt es nicht.

**Stand der Entscheidung.** Auftrag des Nutzers vom 2026-10-18: „Nachweis aus dem Container oder Worktree übernehmen — das hört sich
pragmatisch an“. Die Stelle des Slice (aus den Wellen gelöst, nach dem Kettenumbau) und die Nicht-Änderung der zweiten Sperre des
Hooks (Punkt 20) sind Entscheidungen des Nutzers vom 2026-10-18, so vom Koordinator übermittelt; der Wortlaut liegt dem Architect nicht
vor. Alles Übrige in dieser Kennung ist Entscheidung des Architect und Empfehlung, bis der Nutzer [ADR-0049](../../docs/plan/adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md)
annimmt; ändert er sie, ändert sich nur diese Kennung mit dem Skript.

**Schnittstelle.** Das Fragment `harness/mk/nachweis-uebernehmen.mk` definiert das Ziel `nachweis-uebernehmen` (`.PHONY`) mit einem
Rezept, das weder `$(MAKE)` noch `+` trägt und `bash tools/harness/nachweis-uebernehmen.sh "$(NACHWEIS_QUELLE)"` ruft;
`NACHWEIS_QUELLE` kommt von der Kommandozeile oder aus der Umgebung. Die Quelle ist **das Verzeichnis, das die beiden Dateien
`gates-passed.dateiliste` und `gates-passed.diffsha` trägt**, in einem Worktree also `<worktree>/.harness/state`; die Wurzel
eines Arbeitsbaums ist keine Quelle. Das Skript liest die Quelle einmal, in ein Hilfsverzeichnis unter `${TMPDIR:-/tmp}`, und prüft die
Kopien. Wie die Dateien in ein Verzeichnis kommen (`ssh` mit `tar`, `docker cp`, ein gemeinsamer Mount), ist nicht Teil des Werkzeugs;
das Rezept steht in `docs/developer/devcontainer.md`, ist dort als nicht geprüft gekennzeichnet, und die Gegenprobe führt nur den
Fall mit einer lokalen Kopie (`AGENTS.md` §3.11).

**Plattform** (`AGENTS.md` §3.12, Plattform von [`LH-QA-03`](../../spec/lastenheft.md#lh-qa-03--portabilität)). Zugesagt ist ein Linux-Host
oder -Container mit bash ≥ 4, GNU `make` und GNU-Werkzeugen (`cmp`, `mktemp`, `mv -T`, `readlink -z`, `sort -z`, `sha256sum`), wie die Dateiliste
sie voraussetzt ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 11). macOS- und Windows-Hosts sind nicht zugesagt (akzeptiertes Negativ;
Grund: Die Erzeugung der Liste ist es auch nicht, und das Werkzeug benutzt sie; das Skript lehnt sie nicht ab, und ohne die GNU-Werkzeuge
endet der Lauf in `ungleich` oder mit einem Fehler des Werkzeugs, nie mit `übernommen`).

**Aufbau der Gegenprobe.** `make nachweis-uebernehmen-gegenprobe` (Fragment `harness/mk/nachweis-uebernehmen-gegenprobe.mk`, Skript
`tools/harness/nachweis-uebernehmen-gegenprobe.sh`) steht in `GEGENPROBEN_CONTAINER` und nicht in `GATE_CHECKS`
([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 2). Sie arbeitet in Temp-Repos, ohne Docker und ohne Netz, und lässt den Arbeitsbaum unberührt.
Der Name jedes Falls trägt die Nummer des Punktes. Je Zusage und je Element einer Aufzählung gibt es einen Fall und eine Mutation des
Skripts, die ihn rot färbt; die Gegenprobe verlangt je Mutation genau die genannte Menge roter Fälle (mehr und weniger machen sie rot),
und ohne Nennung ist es der Fall, bei dem die Mutation steht. Wo ein Mechanismus mehrere Elemente deckt, nennt der Punkt die Menge.

**Grundzustand der Fälle.** Zwei Temp-Repos `quelle` und `ziel` desselben Commits mit Kopien von `tools/harness/nachweis-uebernehmen.sh`,
`dateiliste.sh`, `working-tree-hash.sh`, des Hooks und der `.gitignore`; das Verzeichnis `nachweis/` außerhalb beider Bäume trägt Liste und Hash
des Baums `quelle` (vom Rahmen nach der Form von `record-gates.sh` geschrieben, ohne `make gates`). `ziel` trägt in `.harness/state/` einen
älteren gültigen Nachweis eines anderen Baums und einen Stand der schnellen Kette; `TMPDIR` zeigt auf ein leeres Verzeichnis außerhalb. Ein
Fall ruft, wenn nichts anderes steht, das Skript im Baum `ziel` mit `nachweis/` als Quelle und erwartet Exit 0. **Menge A** sind alle Fälle,
die mit Exit ungleich 0 enden; jeder prüft zusätzlich Punkt 14. Die Fälle von Punkt 13 prüfen, dass `TMPDIR` danach leer ist.

**Randform → Punkt.**

| Randform des Plans | Punkt |
|---|---|
| 1 Aufruf ohne Quelle | 3 |
| 2 Nachweis unvollständig | 5 |
| 3 Liste beschädigt | 6 |
| 4 Hash neben der Liste | 7 |
| 5 Sonderzeichen im Pfad der Quelle | 8 |
| 6 Liste ungleich | 9 |
| 7 Gleichheit trotz Sonderpfad | 9 |
| 8 Baum ohne Liste, kein Arbeitsbaum | 11 |
| 9 Index-Eintrag nicht aufgelöst | 12 |
| 10 nicht Teil der Liste | 20 |
| 11 Erfolg ersetzt | 13 |
| 12 Ablehnung ändert nichts | 14 |
| 13 Wiederholung und Alter | 13 |
| 14 Schreiben scheitert oder wird abgebrochen | 15 |
| 15 Quelle ist das Ziel | 4 |
| 16 Quelle im Arbeitsbaum des Ziels | 10 |
| 17 Make-Flags | 16 |
| 18 Aufruf aus einem Unterverzeichnis | 17 |
| 19 gleichzeitige Läufe | 19 |
| 20 Hook liest das Ergebnis | 18 |
| 21 jeder lokale Nachweis ist von Hand zu schreiben | 21 |
| 22 der Nachweis gilt einem Baum | 22 |
| 23 Reihenfolge der Gründe (Architect) | 2 |
| 24 Plattform (Architect) | Kopf dieser Kennung |
| 25 Hilfsverzeichnis und `TMPDIR` (Architect) | 14, 15 |
| 26 Quelle und Stand der schnellen Kette bleiben (Architect) | 14 |

**Ausgang und Meldung.**

1. **Exit-Codes, Kennwörter, Ausgabe.** Exit 0: übernommen; Exit 1: abgelehnt, weil der Nachweis oder der Baum es nicht trägt; Exit 2:
   Fehler im Aufruf. Die erste Zeile auf stderr einer Ablehnung oder eines Aufruffehlers ist `nachweis-uebernehmen: <kennwort>: <text>`,
   zugesagt sind Exit-Code und Kennwort, nicht der Text (außer wo ein Punkt ihn nennt). Die Kennwörter sind genau diese, mit ihrem Exit-Code:
   `make-i` (2), `kein-arbeitsbaum` (2), `aufruf` (2), `quelle-kein-verzeichnis` (2), `quelle-ist-ziel` (2), `kein-nachweis` (1),
   `liste-beschaedigt` (1), `hash-passt-nicht` (1), `ziel-ohne-liste` (1), `ungleich` (1), `schreiben` (1). Der Erfolg schreibt auf stdout
   eine Zeile, die mit `nachweis-uebernehmen: übernommen` beginnt. Über `make` ist nur „Exit ungleich 0“ zugesagt (`make` gibt einem
   scheiternden Rezept Exit 2, bei `-i` Exit 0); die Gegenprobe prüft die Codes am Skript und über `make` das Rezept, die Meldung und das
   Ergebnis. Das Skript schreibt keine weitere Ausgabe auf stdout oder stderr als diese Zeilen. Das Ziel steht in keiner Hälfte von
   `GATE_CHECKS` und nicht in `GATE_SCHNELL_AUSNAHME`; die Gegenprobe steht in `GEGENPROBEN_CONTAINER` und in keiner Hälfte von `GATE_CHECKS`.
   Fall `p1-kein-gate` (Temp-Baum mit Makefile und beiden Fragmenten): `make -n gates` nennt das Ziel nicht, es steht in keiner Kette, und die
   Gegenprobe steht in der Liste; Mutation: das Fragment hängt das Ziel an `GATE_CHECKS`. Die Exit-Codes und Kennwörter prüft jeder Fall der
   Menge A und der Erfolg in Punkt 13; Mutationen: Exit 1 und 2 vertauscht (färbt alle Fälle der Menge A); das Kennwort fehlt in der Meldung
   (alle Fälle der Menge A); Exit 0 trotz Ablehnung (alle der Menge A).
2. **Reihenfolge der Gründe.** Das Skript prüft in dieser Reihenfolge und endet beim ersten Treffer: (1) `-i` in `MAKEFLAGS` (`make-i`),
   (2) Ziel ist kein Git-Arbeitsbaum (`kein-arbeitsbaum`), (3) Aufruf (`aufruf`), (4) Quelle ist kein Verzeichnis (`quelle-kein-verzeichnis`),
   (5) Quelle ist das Ziel (`quelle-ist-ziel`), (6) Datei der Quelle fehlt, ist leer oder kein lesbarer Inhalt (`kein-nachweis`),
   (7) Liste beschädigt (`liste-beschaedigt`), (8) Hash passt nicht (`hash-passt-nicht`), (9) frische Liste des Ziels nicht zu bilden oder leer
   (`ziel-ohne-liste`), (10) Listen ungleich (`ungleich`), (11) Schreiben (`schreiben`). Die Prüfungen (1) bis (10) ändern nichts am Arbeitsbaum.
   Fälle, je Nachbarschaft einer, jeweils mit beiden Gründen zugleich: `p2-i-vor-aufruf` (`-i` und leere Quelle: `make-i`),
   `p2-arbeitsbaum-vor-aufruf` (Aufruf außerhalb eines Git-Arbeitsbaums mit leerer Quelle: `kein-arbeitsbaum`), `p2-ziel-vor-nachweis` (Quelle ist
   `.harness/state` des Ziels ohne Nachweis: `quelle-ist-ziel`), `p2-nachweis-vor-liste` (Liste beschädigt, Hash fehlt: `kein-nachweis`),
   `p2-liste-vor-hash` (Liste beschädigt, Hash falsch: `liste-beschaedigt`), `p2-hash-vor-ohneliste` (Hash falsch, frische Liste leer:
   `hash-passt-nicht`), `p2-hash-vor-ungleich` (Hash falsch, Listen ungleich: `hash-passt-nicht`), `p2-ungleich-vor-schreiben` (Listen
   ungleich, `.harness/state` ist eine Datei: `ungleich`). Mutationen: je Nachbarschaft die beiden Prüfungen vertauscht (färbt den Fall dieser
   Nachbarschaft).

**Teil 1: die Quelle.**

3. **Aufruf.** Ein Aufruf ohne Argument, mit einem leeren Argument (`NACHWEIS_QUELLE` ungesetzt oder leer) und mit mehr als einem Argument endet
   mit Exit 2, `aufruf`. Ein Argument, das nicht auf ein Verzeichnis zeigt (nicht vorhanden, eine Datei), endet mit Exit 2,
   `quelle-kein-verzeichnis`. Ein Symlink auf ein Verzeichnis zählt als Verzeichnis. Fälle: `p3-ungesetzt`, `p3-leer`, `p3-zwei-argumente`,
   `p3-nicht-vorhanden`, `p3-ist-datei`. Mutationen: ohne Argument wird `.` als Quelle gelesen (`p3-ungesetzt`); ein leeres Argument gilt als
   `.` (`p3-leer`); das zweite Argument wird verworfen (`p3-zwei-argumente`); ein nicht vorhandenes Verzeichnis endet mit Exit 1
   (`p3-nicht-vorhanden`); eine Datei als Quelle wird als Liste gelesen (`p3-ist-datei`). Ein Verzeichnis ohne Leserecht ist unter `root` nicht
   herstellbar und hat keinen Fall (akzeptiertes Negativ; Grund: es endet als `kein-nachweis`, und `p5-liste-verzeichnis` deckt den Pfad „Datei nicht
   lesbar“ ohne Rechte).
4. **Quelle ist das Ziel.** Zeigt die Quelle auf dasselbe Verzeichnis wie `.harness/state` des Ziels (Vergleich der Verzeichnisse selbst, `-ef`),
   endet das Skript mit Exit 2, `quelle-ist-ziel`, in jeder Schreibweise: gleicher Pfad, Pfad mit `./`-Segment, Pfad mit Schrägstrich am Ende,
   Symlink auf das Verzeichnis. Fälle: `p4-gleich`, `p4-punkt-segment`, `p4-schraegstrich`, `p4-symlink`. Mutationen: die Prüfung gestrichen (färbt
   alle vier); `-ef` durch einen Textvergleich der Rohpfade ersetzt (färbt `p4-punkt-segment`, `p4-schraegstrich`, `p4-symlink`); der Textvergleich
   wird um Schrägstriche am Ende bereinigt (färbt `p4-punkt-segment`, `p4-symlink`).
5. **Nachweis der Quelle vollständig.** Fehlt `gates-passed.dateiliste` oder `gates-passed.diffsha` im Verzeichnis, ist sie leer oder
   ist sie ein Verzeichnis, endet das Skript mit Exit 1, `kein-nachweis`. Ein Symlink auf eine Datei der Quelle wird gelesen wie die Datei; ein
   Symlink ist darüber hinaus nicht zugesagt. Fälle: `p5-liste-fehlt`, `p5-hash-fehlt`, `p5-liste-leer`, `p5-hash-leer`, `p5-liste-verzeichnis`,
   `p5-hash-verzeichnis`. Mutationen: je Fall die Prüfung dieser Form gestrichen (färbt ihn; in `p5-liste-leer` fängt `liste-beschaedigt` den Fall, das
   Kennwort ist falsch).
6. **Liste unbeschädigt.** Eine Liste, die `dateiliste_pruefen` ablehnt, endet mit Exit 1, `liste-beschaedigt`; der Hash der Fälle passt zur
   beschädigten Liste, damit nur die Beschädigung den Fall bestimmt. Beschädigt ist (Aufzählung wie [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 20,
   leer steht in Punkt 5): Feldzahl nicht durch 3 teilbar, unbekannte Art, `F` mit Wert, der kurz ist, `F` mit Wert aus 64 Zeichen, die nicht
   hexadezimal sind, `G` mit Wert, `X` mit Wert, doppelter Pfad. Fälle: `p6-feldzahl`, `p6-art-unbekannt`, `p6-f-kurz`, `p6-f-nichthex`, `p6-g-wert`,
   `p6-x-wert`, `p6-pfad-doppelt`. Mutationen: je Fall die Prüfung dieser Form aus `dateiliste_pruefen` der Kopie gestrichen (färbt ihn; das
   Kennwort ist dann `ungleich`). Großbuchstaben in einem `F`-Wert besteht `dateiliste_pruefen` (Bestand) und endet in `ungleich`; eine am
   Eintragsende abgeschnittene Liste ist nicht beschädigt und endet ebenfalls in `ungleich`, Punkt 9 (`p9-zu-wenig`) (beide ohne eigenen Fall).
7. **Hash neben der Liste.** `gates-passed.diffsha` ist byte-gleich `<SHA-256 der Liste, 64 Kleinbuchstaben-Hexzeichen>\n`; sonst Exit 1,
   `hash-passt-nicht`. Das fängt auch eine halb geschriebene Quelle (Hash neu, Liste alt). Fälle: `p7-anderer-hash` (der Hash einer anderen Liste),
   `p7-ohne-umbruch`, `p7-grossbuchstaben`, `p7-zweite-zeile`. Mutationen: der Hash wird nicht geprüft (färbt alle vier); ohne Zeilenumbruch
   gilt als gleich (`p7-ohne-umbruch`); der Vergleich ohne Beachtung der Schreibweise (`p7-grossbuchstaben`); nur die erste Zeile
   (`p7-zweite-zeile`); die Prüfung des Hashs einer anderen Liste entfällt, geprüft wird nur die Form (`p7-anderer-hash`).
8. **Sonderzeichen im Pfad der Quelle.** Ein Pfad mit Leerzeichen und ein Pfad, der mit `-` beginnt (relativ übergeben), werden angenommen wie
   jeder andere. Andere Sonderzeichen im Pfad (Zeilenumbruch, Tabulator, Nicht-UTF-8) sind nicht zugesagt. Fälle: `p8-leerzeichen`, `p8-minus`.
   Mutationen: der Pfad ohne Anführungszeichen (`p8-leerzeichen`); ohne `--` oder ohne absoluten Pfad vor dem Lesen (`p8-minus`).

**Teil 2: Gleichheit mit dem Baum, der hier liegt.**

9. **Gleichheit.** Die Liste der Quelle ist byte-gleich mit `dateiliste_erzeugen` im Ziel; sonst Exit 1, `ungleich`, und die Meldung nennt mindestens
   einen abweichenden Pfad (fehlt, neu oder anderer Wert; Zahl und Darstellung darüber hinaus nicht zugesagt, für einen Pfad mit Zeilenumbruch
   auch die Darstellung nicht). Fälle: `p9-hash` (anderer Inhaltshash eines Pfades in der Liste der Quelle), `p9-zu-viel` (ein Pfad mehr in der Quelle), `p9-zu-wenig`
   (ein Pfad weniger), `p9-art` (`F` gegen `L`), `p9-link-ziel`, `p9-untracked` (eine untracked Datei im Ziel), `p9-geaendert` (eine getrackte Datei im Ziel
   geändert), `p9-geloescht` (eine getrackte Datei im Ziel gelöscht, Art `G`). Zwei Fälle mit Erfolg: `p9-sonder-umbruch` (ein Pfad mit
   Zeilenumbruch in beiden Listen) und `p9-sonder-leerzeichen` (ein Pfad mit Leerzeichen): übernommen. Mutationen: das Skript bildet die Liste selbst
   (`git ls-files` ohne `--others`; färbt `p9-untracked`); es vergleicht nur die Pfade (färbt `p9-hash`, `p9-art`, `p9-link-ziel`, `p9-geaendert`,
   `p9-geloescht`); es verlangt nur, dass jeder Pfad des Ziels in der Quelle steht (`p9-zu-viel`); es verlangt nur, dass jeder Pfad der Quelle im Ziel steht
   (`p9-zu-wenig`); es liest die Liste zeilenweise statt NUL-getrennt (färbt `p9-sonder-umbruch`); es teilt an Leerzeichen (`p9-sonder-leerzeichen`). Die Meldung
   nennt den Pfad: Mutation, die ihn nicht nennt, färbt `p9-hash`, `p9-zu-viel`, `p9-zu-wenig`, `p9-untracked`, `p9-geaendert`, `p9-geloescht`.
10. **Quelle im Arbeitsbaum des Ziels.** Liegt die Quelle als untracked, nicht ignoriertes Verzeichnis im Baum, sind ihre Dateien Teil der frischen Liste und der
    Vergleich endet in `ungleich`; liegt sie unter `.harness/state/` (ignoriert, aber nicht dieses Verzeichnis selbst), wird sie angenommen. Fälle:
    `p10-untracked-im-baum`, `p10-unter-state`. Mutationen: das Skript schließt die Quelle aus der frischen Liste aus (`p10-untracked-im-baum`); es
    verlangt, dass die Quelle außerhalb des Baums liegt (`p10-unter-state`).
11. **Baum ohne Liste, kein Arbeitsbaum.** Scheitert `git` im Ziel oder ist die frische Liste leer, endet das Skript mit Exit 1, `ziel-ohne-liste`; ist das
    Verzeichnis des Aufrufs kein Git-Arbeitsbaum (auch das Innere von `.git`), endet es mit Exit 2, `kein-arbeitsbaum`. Fälle: `p11-git-scheitert` (Attrappe
    von `git`, `ls-files` endet mit 1), `p11-leer` (alle Dateien des Ziels sind ignoriert, auch das Skript), `p11-kein-arbeitsbaum`. Mutationen: ein
    scheiterndes `git` zählt als leere Liste ohne Fehler (`p11-git-scheitert`: dann `ungleich`); die leere Liste wird mit der Liste der Quelle verglichen
    (`p11-leer`: dann `ungleich`); `kein-arbeitsbaum` endet mit Exit 1 (`p11-kein-arbeitsbaum`).
12. **Nicht aufgelöster Index-Eintrag im Ziel.** Das Skript übernimmt bei gleicher Liste auch bei einem Konflikt (`git ls-files -u` nicht leer), die Liste trägt
    keine Index-Stufen; der Hook sperrt weiter ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 17). Die Meldung sagt dazu nichts zu. Fall `p12-konflikt`: Exit 0, die
    Dateien sind übernommen, und der unveränderte Hook antwortet `block`. Mutation: das Skript lehnt bei einem Konflikt ab.

**Teil 3: das Ziel.**

13. **Erfolg ersetzt; Wiederholung.** Nach `übernommen` sind `gates-passed.dateiliste` und `gates-passed.diffsha` des Ziels byte-gleich denen der Quelle, und das Skript
    hat ihre Reihenfolge eingehalten: zuerst der Hash, dann die Liste (die Datei, die der Hook liest, ist die letzte). Ein älterer Nachweis (gültig, beschädigt), ein fehlender
    und ein fehlendes Verzeichnis `.harness/state/` ändern das nicht (das Verzeichnis wird angelegt). Ein zweites Übernehmen derselben Quelle endet mit Exit 0 und gleichen
    Bytes. Ein Nachweis trägt keinen Zeitpunkt; ein älterer Nachweis desselben Baums ist so gültig wie ein neuerer. Dateimodi sind nicht zugesagt. Fälle: `p13-neu` (kein
    älterer Nachweis; prüft auch die Zeile `nachweis-uebernehmen: übernommen`), `p13-kein-state`, `p13-aelter-gueltig`, `p13-aelter-beschaedigt`, `p13-zweimal`.
    Mutationen: das Skript schreibt nur die Liste (`p13-neu`, alle); es legt `.harness/state/` nicht an (`p13-kein-state`); es lehnt ab, wenn ein Nachweis da ist
    (`p13-aelter-gueltig`, `p13-aelter-beschaedigt`, `p13-zweimal`); es lehnt ab, wenn der Nachweis gleich ist (`p13-zweimal`); es ändert die Bytes (fügt einen Zeilenumbruch an;
    färbt alle Fälle dieses Punktes).
14. **Ablehnung ändert nichts; Quelle und Stand der schnellen Kette bleiben.** Jeder Fall der Menge A prüft nach dem Lauf: die beiden Dateien in `.harness/state/` des Ziels sind
    byte-gleich dem Zustand davor, der Stand der schnellen Kette ebenso, `.harness/state/` trägt keine neue Datei (keine Hilfsdatei), `TMPDIR` ist leer, und die Quelle ist
    byte-gleich dem Zustand davor. Eine Ablehnung legt `.harness/state/` nicht an. Auch im Erfolg bleiben die Quelle und
    der Stand der schnellen Kette (`.harness/state/gates-schnell-passed.diffsha`) byte-gleich; das Skript übernimmt den Stand der schnellen Kette nicht und löscht ihn nicht
    (Entscheidung des Nutzers vom 2026-10-18, so übermittelt: die zweite Sperre des Hooks bleibt, [ADR-0049](../../docs/plan/adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md)). Fälle:
    `p14-ohne-state` (Ablehnung `ungleich` bei fehlendem `.harness/state/`: es fehlt danach noch), `p14-quelle-erfolg`, `p14-quelle-ablehnung`, `p14-schnell-erfolg`, `p14-schnell-ablehnung`.
    Mutationen: eine Ablehnung entfernt den älteren Nachweis (färbt alle Fälle der Menge A); eine Ablehnung lässt eine Hilfsdatei in `.harness/state/` (Menge A); ein Hilfsverzeichnis
    unter `TMPDIR` bleibt nach Ablehnung (Menge A) und nach Erfolg (Fälle von Punkt 13); eine Ablehnung legt `.harness/state/` an (`p14-ohne-state`); die Quelle wird verschoben statt kopiert
    (`p14-quelle-erfolg`); die Quelle wird bei Ablehnung gelöscht (`p14-quelle-ablehnung`); der Erfolg löscht den Stand der schnellen Kette (`p14-schnell-erfolg`); eine Ablehnung löscht ihn
    (`p14-schnell-ablehnung`).
15. **Schreiben scheitert oder wird abgebrochen.** Ist `.harness/state` eine Datei, oder ist eine der beiden Zieldateien ein Verzeichnis, endet das Skript mit Exit 1, `schreiben`, vor der ersten
    Ersetzung; der ältere Nachweis bleibt unverändert, und es bleibt keine Hilfsdatei (Prüfung wie Menge A). `SIGINT` und `SIGTERM` beenden das Skript mit 130 und 143, ohne Hilfsdatei und
    ohne Hilfsverzeichnis; jede der beiden Dateien ist danach die alte oder die neue ganz, nie halb. Ob nach einem Signal zwischen den beiden Ersetzungen das Paar zusammengehört, ist nicht
    zugesagt (akzeptiertes Negativ; Grund: Das Fenster sind zwei `mv`; die Folge, Hash neu und Liste alt, lehnt Punkt 7 ab, wenn jemand das Verzeichnis als Quelle nimmt, und der Hook liest nur die Liste).
    Ein Verzeichnis `.harness/state` ohne Schreibrecht ist unter `root` nicht herstellbar und hat keinen Fall (akzeptiertes Negativ; Grund: Die Fälle mit einer Datei und mit einem Verzeichnis zeigen
    dieselbe Ablehnung ohne Rechte). Fälle: `p15-state-ist-datei`, `p15-liste-ist-verzeichnis`, `p15-hash-ist-verzeichnis`; je Signal und Phase einer, vier insgesamt: `p15-sigint-pruefen`,
    `p15-sigterm-pruefen` (eine Attrappe von `cmp` wartet, das Signal geht an das Skript), `p15-sigint-ersetzen`, `p15-sigterm-ersetzen` (eine Attrappe von `mv` wartet vor der ersten Ersetzung); geprüft
    werden Exit-Code, Zustand wie in Punkt 14 und dass die Attrappe angelaufen war. Mutationen: das Skript prüft die Ziele nicht vor der ersten Ersetzung (färbt `p15-liste-ist-verzeichnis`: Hash neu,
    Liste alt); ein Verzeichnis als Ziel der Liste wird mit `mv` ohne `-T` in das Verzeichnis geschoben (`p15-liste-ist-verzeichnis`); dasselbe für den Hash (`p15-hash-ist-verzeichnis`);
    `mkdir -p` scheitert ohne Fehler (`p15-state-ist-datei`); die Behandlung von `SIGINT` gestrichen (färbt beide Fälle mit `SIGINT`); die von `SIGTERM` gestrichen (beide mit `SIGTERM`); `SIGINT` endet mit 143 und
    `SIGTERM` mit 130 (je zwei Fälle); das Hilfsverzeichnis bleibt beim Signal (`p15-sigint-pruefen`, `p15-sigterm-pruefen`); die Hilfsdatei in `.harness/state/` bleibt beim Signal
    (`p15-sigint-ersetzen`, `p15-sigterm-ersetzen`).

**Teil 4: Aufrufarten.**

16. **Make-Flags.** Zugesagt ist das Ergebnis, nicht der Mechanismus; maßgeblich ist `MAKEFLAGS`, auch aus der Umgebung. (a) `-n`, `-q`, `-t`: `make` führt das Rezept nicht aus (es trägt weder
    `$(MAKE)` noch `+`), nichts ändert sich, auch dann nicht, wenn die Übernahme gelänge; der Exit-Code unter `-q` ist nicht zugesagt. (b) `-i`: Das Skript prüft die Buchstaben im ersten Wort von
    `MAKEFLAGS`, wenn es nicht mit `-` beginnt (wie [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 21), und endet bei `i` mit Exit 2, `make-i`; es ändert nichts. Grund: Die Einzelschritte erbten `-i`,
    und eine Ablehnung sähe unter `make` wie Erfolg aus (Exit 0); die Meldung ist dann das Zeichen. Anders als bei `make gates` entfällt ein vorhandener Nachweis nicht (Grund: Das Werkzeug fährt
    keine Gates; der Nachweis beschreibt weiter den Baum seines grünen Laufs). `-i` gilt in jeder Schreibweise, aus der Umgebung und zusammen mit `-k` oder `-j`. (c) `-i` zusammen mit `-n`, `-q` oder `-t`:
    (a) gilt. (d) `-k`, `-j` und Langflags (`--no-print-directory`) ändern nichts; die Auswahl liest nur das erste Wort. Fälle, am Skript (`MAKEFLAGS` gesetzt) und über `make`: `p16-n`, `p16-q`,
    `p16-t` (Kommandozeile), `p16-n-env`, `p16-q-env`, `p16-t-env` (`MAKEFLAGS` der Umgebung): nichts geändert; `p16-i`, `p16-i-lang` (`--ignore-errors`), `p16-i-env`, `p16-i-k`, `p16-i-j`: am Skript Exit 2
    mit `make-i`, über `make` die Meldung und nichts geändert; `p16-ni`, `p16-qi`, `p16-ti`: nichts geändert, keine Meldung des Skripts; `p16-k`, `p16-j`: übernommen; `p16-langflag` (`MAKEFLAGS=k --no-print-directory`):
    übernommen. Mutationen: das Rezept trägt `+` (färbt `p16-n`, `p16-q`, `p16-t` und die drei `-env`-Fälle); das Skript überliest `i` (färbt `p16-i`, `p16-i-lang`, `p16-i-env`, `p16-i-k`, `p16-i-j`); das Rezept trägt `+` und das Skript überliest `i` (färbt `p16-ni`, `p16-qi`, `p16-ti`, dazu `p16-n`, `p16-q`, `p16-t`, die drei `-env`-Fälle und die fünf `-i`-Fälle); das Skript liest das ganze `MAKEFLAGS` statt des ersten Worts (`p16-langflag`); ohne `k` lehnt das
    Skript ab (`p16-k`); `j` lehnt ab (`p16-j`); die Ablehnung entfernt einen vorhandenen Nachweis (färbt `p16-i`, `p16-i-lang`, `p16-i-env`, `p16-i-k`, `p16-i-j`).
17. **Aufruf aus einem Unterverzeichnis.** Das Skript urteilt über den Baum der Wurzel des Repos, in dem sein Arbeitsverzeichnis liegt (`git rev-parse --show-toplevel`); eine relative Quelle gilt gegen das Arbeitsverzeichnis des
    Aufrufs. Ein Aufruf aus einem Unterverzeichnis über `make -C` ist Sache von `make` (akzeptiertes Negativ; Grund: wie [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 13, das Rezept läuft im Verzeichnis, das `make` wählt). Fälle:
    `p17-relativ`, `p17-absolut`. Mutationen: eine relative Quelle gilt gegen die Wurzel (`p17-relativ`); das Skript wechselt nicht in die Wurzel (`p17-absolut`: die Liste des Unterverzeichnisses ist verschieden).

**Teil 5: Kopplung an den Hook.**

18. **Der unveränderte Hook liest das Ergebnis.** In `ziel` antwortet der Hook aus dem Arbeitsbaum, ohne Änderung: nach `übernommen` `approve`, ohne Stand der schnellen Kette
    ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 16); vor der Übernahme bei fehlendem Nachweis `block` (Punkt 19); nach einer Ablehnung wegen einer untracked Datei `block`; nach einer Ablehnung mit einem
    älteren gültigen Nachweis eines anderen Baums bleibt es bei dessen Befund (`block`, weil der Unterschied nicht inert ist). Fälle: `p18-erfolg` (vorher `block`, nachher `approve`), `p18-ohne-nachweis`, `p18-untracked`,
    `p18-alt`. Der Fall `p18-untracked` legt die untracked Datei außerhalb der Positivliste an. Mutationen: das Skript vergleicht nicht und schreibt die frische Liste des Ziels (färbt die acht Fälle
    der Ablehnung `ungleich` in Punkt 9 und `p18-untracked`, `p18-alt` mit `approve`); das Skript schreibt bei jeder Ablehnung die frische Liste des Ziels (färbt `p18-ohne-nachweis`, `p18-untracked`, `p18-alt` mit `approve`
    und alle Fälle der Menge A); das Skript schreibt keine Liste (`p18-erfolg`: bleibt `block`). Der Hook ist ein zweites Netz: Bildete das Skript die Liste ohne `--others`, träfe der Hook die untracked Datei weiter und
    sperrte; diese Mutation färbt `p9-untracked`, nicht `p18-untracked`.

**Für alle Teile.**

19. **Aufruf-Voraussetzung.** Zugesagt ist nacheinander: eine Übernahme, ein Lauf von `make gates`, eine Änderung des Arbeitsbaums und der Hook kommen nicht gleichzeitig. Gleichzeitige Läufe sind nicht zugesagt und
    haben keinen Fall. Ändert sich der Baum nach der Übernahme, vergleicht der Hook neu; das Werkzeug kann nie öfter freigeben als der Hook, weil dieser den Baum liest. Das Skript sagt keine Ordnung über mehrere Aufrufe zu.
20. **Nicht Teil der Liste, damit nicht geprüft** (Bestand, [`SPEC-052`](SPEC-052-gate-kette.md) Punkte 11 und 18; ohne Fall). Dateimodi, ignorierte Dateien (`.gitignore`, `.git/info/exclude`, `core.excludesFile`), Zeitstempel und `.harness/state/`
    gehören nicht zur Liste. Zwei Umgebungen desselben Commits unterscheiden sich in der Liste durch untracked, nicht ignorierte Dateien, durch Inhalte, die `git` beim Auschecken umschreibt (`core.autocrlf`), und durch lokale
    Ignore-Regeln; jede dieser Abweichungen endet in `ungleich`, nie in einer falschen Übernahme. **Die zweite Sperre des Hooks bleibt** (Entscheidung des Nutzers vom 2026-10-18, so übermittelt): Eine inerte Änderung im
    Hauptbaum nach einer Übernahme verlangt weiter den Stand der schnellen Kette ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 9); das Werkzeug übernimmt ihn nicht (Punkt 14). Eine Übernahme bei ungleicher Liste, deren Unterschiede alle
    auf der Positivliste liegen, gibt es nicht (Grund: Die Positivliste steht an einer Stelle, im Hook, [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 12).
21. **Grenze, offen genannt: jeder lokale Nachweis ist von Hand zu schreiben.** `dateiliste_erzeugen > .harness/state/gates-passed.dateiliste`, `make record-gates` einzeln und `cp` eines fremden Nachweises schreiben ihn heute schon.
    Das Werkzeug ist der einzige vorgesehene Weg über eine Umgebungsgrenze und prüft Gleichheit, Unversehrtheit und Zusammengehörigkeit; es kann nicht prüfen, dass die Quelle `make gates` fuhr, mit denselben Gates, im selben
    Stand des Images. Akzeptiertes Negativ; Grund: Der Nachweis trägt keine Herkunft, und eine Herkunft ohne Signatur wäre von Hand ebenso zu schreiben. Ein Gate, das nur in einer anderen Umgebung rot wird, bleibt in der Quelle
    unbemerkt (`BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung`). Eine Gegenprobe gibt es nicht, weil es nichts gibt, was sie rot färben könnte.
22. **Der Nachweis gilt einem Baum, nicht einem Zweig.** Ein Merge, der den Baum des Hauptzweigs ändert, macht die Listen ungleich. Die Regel an den Aufrufer: den Zweig vor dem Merge auf den Stand des Hauptzweigs bringen und dort
    `make gates` fahren, oder `make gates` im Hauptbaum. Das Werkzeug lehnt ab (`ungleich`), es rät nicht. Ohne Fall.
