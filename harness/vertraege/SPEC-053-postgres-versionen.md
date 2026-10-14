# SPEC-053 — Versionsmatrix der Integrationstests (`test-postgres-versionen`)

`make test-postgres-versionen` fährt die Teilmenge der Integrationstests, deren Aussage vom
Server abhängt, gegen jede Version einer Liste von PostgreSQL-Images außer der Referenz;
`make test-postgres-versionen-gegenprobe` prüft das Werkzeug. Die Referenz fährt
`make test-integration` mit der vollen Suite, unverändert. Diese Kennung legt fest, was das
Gate prüft und woran es rot wird. Welche Versionen das Produkt zusagt, steht nicht hier,
sondern in der Produktspezifikation; dieser Vertrag sagt nur, was das Gate belegt.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Die Gegenprobe
fährt in Temp-Bäumen mit einer Attrappe für `docker` (über `DOCKER`), ohne Docker und ohne Netz;
die unveränderte Kopie ist grün, und je Fall bricht eine Mutation genau ihn.

1. **Liste.** Die Liste steht an einer Stelle, im Fragment `harness/mk/postgres-versionen.mk`.
   Ein Eintrag ist `<Hauptversion>=<Image>` mit `<Image>` als `<Name>:<Tag>@sha256:<Digest>`;
   das Fragment benennt die Hauptversion der Referenz. Die Liste ist die bestätigte Liste zum
   Zeitpunkt der Lieferung; die Bestätigung gegen die Release-Liste von PostgreSQL ist eine
   Handlung des Slice, kein Lauf des Gates (Punkt 14).
2. **Eingabe, rot vor dem ersten Aufruf von Docker.** Das Ziel endet mit Exit-Code ungleich 0,
   nennt den Eintrag und ruft `docker` nicht auf bei: (a) leerer Liste, (b) Referenz fehlt
   in der Liste, (c) Hauptversion kommt zweimal vor (auch die Referenz; ein Alias auf dieselbe
   Hauptversion), (d) keine andere Version als die Referenz, (e) Eintrag ohne `=`, mit leerer
   oder nicht aus Ziffern bestehender Hauptversion oder mit führender Null, (f) Image ohne Tag,
   (g) Image ohne Digest, (h) Digest, der nicht aus 64 kleinen Hexziffern besteht, (i) Image
   der Referenz ungleich `POSTGRES_IMAGE` aus `harness/mk/integration.mk`. Jede Form hat einen
   Fall und eine Mutation, die nur sie bricht.
3. **Teilmenge.** Die Teilmenge steht im Fragment an einer Stelle: Phase 1 (mit PostgreSQL)
   und Phase 2 (ohne PostgreSQL, gegen die Aufzeichnungen der Phase 1), je mit dem Namen jedes
   Tests, der laufen muss. Sie ist die Abnahmeszenarien 1, 2 und 7, der Fall SCRAM-SHA-256 der
   Anmeldung und der Abgleich der Lebendprüfung mit dem Server. Ein erwarteter Test, dessen
   Ausgabe keine `PASS`-Zeile trägt (umbenannt, nicht gewählt, übersprungen), macht die Version
   rot; die Meldung nennt Test, Version und Phase, und die Phase ist die, in der der Test
   erwartet war: fehlt er in Phase 2, heißt sie „Phase 2“. Eine Auswahl ohne Treffer ist damit rot,
   nicht grün. Je Form ein Fall und eine Mutation: ein umbenannter Test, ein übersprungener Test (`SKIP`), eine Auswahl ohne Treffer;
   dazu ein Fall mit einem in Phase 2 fehlenden Test, der „Phase 2“ nennt, und eine Mutation, die die Phase der Prüfung festlegt.
4. **Server.** Das Ziel liest nach dem Start die Hauptversion des laufenden Servers und
   vergleicht sie mit der des Eintrags. Rot, mit Nennung der Version, sind: ein Server, der
   nach 60 Sekunden nicht bereit ist; ein Image, das sich nicht starten lässt (der Start zieht
   es, wenn es fehlt; Ziehen und Starten sind ein Schritt und haben einen Fall); eine
   Hauptversion, die nicht lesbar ist oder von der des Eintrags abweicht (Digest, der zu einem
   anderen Tag gehört). Ein Vergleich nur am Tag des Images gilt nicht. Jede dieser Formen und ein Fehler der Einrichtung
   der Anmeldung (Punkt 5) endet vor Phase 1 und steht in der Zusammenfassung als `<Version> (Aufbau)`
   (Punkt 6); jeder Fall dieser Formen prüft diese Stufe.
5. **Anmeldung.** Je Version legt der Lauf einen Benutzer mit Passwort an und trägt vor der
   allgemeinen Zeile in `pg_hba.conf` die Methode `scram-sha-256` für ihn ein. Zugesagt ist der
   Fall SCRAM-SHA-256; Klartext-Passwort und MD5 laufen in der Matrix nicht, und wie ihre
   Einrichtung auf PostgreSQL 18 verläuft, ist nicht zugesagt. Scheitert die Einrichtung, ist
   die Version rot.
6. **Ablauf und Meldung.** Die Versionen laufen nacheinander in der Reihenfolge der Liste, nie
   zwei Server gleichzeitig; die Referenz fährt das Ziel nicht. Eine Version ist rot in einer von drei
   Stufen: „Aufbau“ (vor Phase 1: Punkt 4 und 5), „Phase 1“ und „Phase 2“; die Meldung und die
   Zusammenfassung nennen die Stufe. Ist Phase 1 einer Version rot, entfällt ihre Phase 2. Das Ziel fährt alle Versionen auch
   nach einer roten, fasst am Ende die roten Versionen mit Phase zusammen und endet mit Exit-Code
   ungleich 0, wenn eine rot war. Ein Fehler der Eingabe (Punkt 2) endet vor der ersten Version.
   Das Testimage wird einmal gebaut, nicht je Version. Ist der Bau rot, endet das Ziel sofort mit
   Exit-Code ungleich 0 und einer Meldung, dass das Testimage nicht zu bauen ist; es startet keine
   Version (kein Netz, kein Server) und fasst keine Version zusammen. Der Fall prüft Exit-Code, Meldung
   und dass kein Aufruf von `docker` außer dem Bau stattfand.
7. **Aufräumen.** Container, Netz und Volume einer Version sind entfernt, bevor die nächste
   startet, und nach Ende des Ziels, auch nach einem roten Lauf, einem Fehler des Aufbaus und
   nach `SIGINT` und `SIGTERM`; je Fall des Endes ein Fall. Bei `SIGINT` und `SIGTERM` gilt die Zusage für das
   Skript `tools/test/run-postgres-versionen.sh`: Es räumt auf, bevor es endet, ob das Signal an das
   Skript oder an die Prozessgruppe geht; je Signal und Ziel ein Fall. Ob nach dem Ende von `make`
   nichts mehr da ist, ist nicht zugesagt (Punkt 14). `SIGKILL` und ein Ausfall des
   Docker-Dienstes sind nicht zugesagt (Punkt 14). Die Namen tragen die Kennung des Laufs und
   die Hauptversion.
8. **Der Lauf der Referenz bleibt.** Der Runner `tools/test/run-integration-tests.sh` fährt ohne
   Auswahl wie vor diesem Vertrag: volle Suite in Phase 1, in Phase 2 die Auswahl
   `^TestE2EOhnePostgres`, gleicher Exit-Code, gleiche Tests. Die Ausgabe der Tests auf stderr
   erscheint auf stdout, weil der Runner beide für die Prüfung der `PASS`-Zeilen zusammenführt; die
   Trennung der Kanäle ist nicht zugesagt (Punkt 14). Eine Auswahl der Teilmenge gilt nur, wenn das Ziel sie übergibt.
9. **Aufrufarten des Treibers.** Das Ziel schreibt keinen Stand; geglaubt wird allein sein
   Exit-Code in der Kette von `make gates` ([SPEC-052](SPEC-052-gate-kette.md) Punkt 1). Das
   Rezept besteht aus einer Zeile ohne `$(MAKE)` und ohne `+`. Damit fährt es unter `-n`, `-q`
   und `-t` kein Docker und schreibt nichts, und `MAKEFLAGS` aus der Umgebung gilt wie das
   Flag. Unter `-i` überliest `make` den Exit-Code des Skripts: Ein Aufruf mit `-i` belegt
   nichts, und `make -i gates` schreibt keinen Nachweis (SPEC-052 Punkt 6 und 21). Bei Kombination
   von `-n`, `-q` oder `-t` mit `-i` gilt `-n`, `-q`, `-t`: kein Docker. Je Flag `-n`, `-q`,
   `-t`, `-i`, je `MAKEFLAGS` aus der Umgebung und für `-n` zusammen mit `-i` ein Fall; die
   Mutation ist ein Rezept mit `$(MAKE)`.
10. **Aufruf-Voraussetzung.** Zugesagt ist nacheinander, wie in SPEC-052 Punkt 8. `make -j gates`
    fährt `test-integration` und dieses Ziel gleichzeitig; das ist nicht zugesagt und hat keinen
    Fall (Punkt 14).
11. **Versionsabhängige Aussage.** `TestE2EReplayLebendpruefungWiePostgres` liest die
    Hauptversion aus `server_version` des Servers (führende Ziffern) und prüft die Texte mit
    `\v` gegen die Grenze der Produktspezifikation, am Server und am Replay. Ein Wert ohne
    führende Ziffer macht den Test rot. Weicht eine Version von der Grenze ab, ist das Ziel rot;
    die Berichtigung der Grenze gehört nicht zu diesem Vertrag, sie ist ein eigener Slice mit
    Berichtigung der Produktspezifikation.
12. **Gegenprobe.** `test-postgres-versionen-gegenprobe` hängt an `GATE_CHECKS`. Die
    Attrappe von `docker` protokolliert ihre Aufrufe, liefert Version und Testausgabe nach
    Vorgabe des Falls und lässt Punkt 2 bis 9 prüfen: kein Aufruf bei Eingabefehlern, genau ein
    `build`, die Reihenfolge, die Zusammenfassung roter Versionen, die Entfernung von
    Container, Netz und Volume.
13. **Ausnahmeliste der schnellen Kette.** `test-postgres-versionen` steht in der Ausnahmeliste
    (SPEC-052 Punkt 2 und 4): Es startet PostgreSQL-Container wie `test-integration`, und die
    schnelle Kette fährt solche Gates nicht. Seine Gegenprobe steht dort nach ihrem Namen. Das
    Ziel liest nichts unter `docs/plan/planning/`, `docs/reviews/`, `docs/user/` und nicht
    die `README.md` (SPEC-052 Punkt 13): Der Build-Kontext schließt sie aus, und die Gegenprobe
    legt ihre Bäume selbst an.
14. **Grenzen.** (1) Die Release-Liste von PostgreSQL liest kein Lauf; das Gate prüft weder,
    dass die Liste die fünf jüngsten erschienenen Hauptversionen trägt, noch deren Zahl (Netz
    zur Laufzeit; Grund: Pflege ist eine Handlung). (2) `SIGKILL`, ein Ausfall des
    Docker-Dienstes und ein Zeitlimit für den ganzen Lauf sind nicht zugesagt; einzelne Tests
    tragen ihre eigenen Zeitlimits. (3) Gleichzeitige Läufe, auch `make -j gates`, sind nicht
    zugesagt. (4) Ob ein Test der Teilmenge inhaltlich eine Aussage über die Serverversion
    trägt, prüft das Gate nicht; es prüft, dass er läuft und besteht. (5) Abnahmeszenario 18 und
    Szenario 16 Satz 3 laufen nur gegen die Referenz. (6) Die Protokollversion 3.2 prüft die
    Matrix nicht; ein Client mit 3.0 gegen jeden Server der Liste ist ihr Gegenstand. (7) Über
    `make`: Bei `SIGTERM` an die Prozessgruppe endet `make` vor dem Skript, und dessen Aufräumen folgt
    bis zu einige Sekunden danach; ein `SIGINT` nur an `make` lässt es das Skript abwarten. Beides
    ist Verhalten von `make`, das Rezept kann es nicht ändern, und ein Fall auf den Zeitpunkt nach
    dem Ende von `make` wäre zeitabhängig; kein Lauf der Kette (`make gates`) startet danach ein
    Ziel, das die Reste stört. (8) Die Kanäle stdout und stderr der Tests sind im Lauf des Runners
    nicht getrennt (Punkt 8).
