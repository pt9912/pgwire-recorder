# Verifikation: slice-v1-abschluss-einspielen-tls — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-einspielen-tls.md`, F-601 bis F-608); die Übergaben prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `3d0bd0c..a9c9330` (32 Dateien, +2775 −136). Architect: `0cb6c4f`, `c84df26`, `3a2f1f3`. Implementer: `147afea`, `24bc173`; Nacharbeit `058c677`, `a9c9330`. Review: `b221a16`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `a9c9330`, besonders §1 (*Übernimmt*, *Ausdrücklich NICHT*), §2 (DoD), §3 (mit den Aufträgen aus der Review), §6 (*Randformen*, *Risiken*) und §7 (*Belege des Implementers*, Mutationstabelle, Namensabgleich)
- `spec/spezifikation.md` `LH-FA-20.a` *Start*, *TLS*, *Anmeldung*, *Abbruch im Aufbau*, *Fehlerregeln*; `LH-FA-17.a` (*Wirkung einer URL*, *Fehler*); `spec/architecture.md` §4.4 und §6
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md) ganz, [ADR-0019](../plan/adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) (*Entscheidung*), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) und [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) (*Entscheidung*); [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) über den Lauf von `make a-check`; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)
- am Stand `a9c9330`: `einspielen_tls.go`, der Diff von `einspielen.go`, `zertifizierungsstelle.go`, die Diffs von `cli.go`, `upstream.go`, `verbindung.go`, `bootstrap.go`, `testhilfen_import_test.go`, `proxy_test.go`, `play_tls_e2e_test.go` (F-603)
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- `AGENTS.md` §3.9 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `a9c9330`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch-Präfix `ver-tls-`, nie im Repo.

- **Binary:** `git archive HEAD` in eine frische Kopie, daraus die Stufe `runtime` des `Dockerfile` (distroless) mit eigenem Tag. Jeder `play`-Lauf ist ein eigener Container des Binaries aus dieser Stufe; die Wirkung las ich mit `psql` als Administrator im Server.
- **Unabhängige Gegenstellen (Aufgabe 2):** ein internes Docker-Netz mit festen Adressen und
  - **vier echten PostgreSQL-17-Servern** (das gepinnte Image aus `harness/mk/integration.mk`): `pgssl` und `pgssl2` mit `ssl=on` und demselben Zertifikat, `pgexp` mit einem abgelaufenen Zertifikat, `pgplain` mit `ssl=off`. Die Zertifikate hat `openssl` 3.5 in einem Hilfscontainer (alpine) erzeugt: zwei CAs (`ca.pem`, `ca2.pem`), ein Blatt mit `DNS:pgssl, DNS:localhost, IP:172.31.77.10`, eines für `pgexp` gültig vom 2020-01-01 bis 2020-01-02. `pg_hba.conf` der SSL-Server: `hostnossl … reject`, danach `hostssl` mit `scram-sha-256`, `md5`, `password` für drei Rollen und `trust` für den Rest. Der Server schreibt `log_connections`.
  - **einem Fake-Server** (Python, `socket`), der auf das `SSLRequest` mit `N`, einem anderen Byte, `S` und Rauschen, Ende ohne Antwort oder `S` und Schweigen antwortet und alles protokolliert, was danach noch eintrifft.
  - Die Aufzeichnung (`CREATE TABLE`, `INSERT`, `SELECT`) entstand mit `record` gegen `pgplain`. Der Proxy der Tests (`tlsproxy`) kam in keinem dieser Läufe vor.
- **Unit-Mutanten:** je Mutant eine frische Kopie des Quellstands per `cp -r` ohne `-p`, genau eine Ersetzung (das Skript bricht ab, wenn der Suchtext nicht genau einmal trifft), `gofmt -l` auf der Datei (sauber), im Image der Stufe `deps` mit der Kopie als Bind-Mount, `--network none`: `go test -count=1` über `./internal/adapters/driven/postgres`, `./internal/adapters/driving/cli` und `./internal/bootstrap/...`. Der Grundlauf ohne Mutation war grün.
- **Integrations-Mutant:** eine Kopie, `tools/test/run-integration-tests.sh` mit umbenanntem Image und Lauf-Präfix.
- **`-race`:** `golang:1.27.2` (mit `gcc`), Quellbaum schreibgeschützt, Modulcache als Volume aus dem `deps`-Image, `--network none`.
- **Danach:** Container, Netz, Volumes, Kopien und eigene Images entfernt (siehe Ende). `make gates` lief im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `play` baut mit `--upstream-tls` oder `sslmode=require` TLS auf, prüft das Zertifikat gegen den Systemspeicher plus `--upstream-ca` und den Namen gegen den eingesetzten Host; ein gesetztes `--upstream-tls`, auch `false`, geht vor; ohne beides kein TLS. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| TLS gegen einen echten Server mit eigener CA, Wirkung in der Tabelle | PostgreSQL 17 `ssl=on`, Verbindung aus der Konfigurationsdatei mit `sslmode=require` und `--upstream-ca ca.pem`: Exit 0, die Tabelle trägt `tls-ok`. Der Server-Log zeigt für die Rollen von `play` `connection authorized … SSL enabled (protocol=TLSv1.3 …)`. Ebenso mit `--upstream pgssl:5432 --upstream-tls` (host:port) und mit dem Host als IP-Adresse (`172.31.77.10`, steht im Zertifikat unter den IP-Adressen) | bestätigt |
| Unbekannte CA ist `PGR-E4005`, kein Klartext | ohne `--upstream-ca`: Exit 4, `PGR-E4005`, „certificate signed by unknown authority“; mit der fremden `ca2.pem`: ebenso. Die Tabelle fehlt in beiden Fällen. Das Zertifikat wird also zuerst geprüft, bevor ein Startup geht (der TLS-Aufbau scheitert, bevor Daten fließen) | bestätigt |
| Systemspeicher: `SSL_CERT_FILE` ohne `--upstream-ca` | Binary der Stufe `runtime`, `SSL_CERT_FILE=/work/ca.pem`: Exit 0, Tabelle gefüllt; ohne die Variable `PGR-E4005`. Der Systemspeicher ist also der Standardpfad (Unit-Mutant V05, E2E-Mutant siehe Punkt 4) | bestätigt |
| Ergänzen, nicht ersetzen | `--upstream-ca` mit Datei `Text davor / ca2.pem / Zwischentext / ca.pem / danach`: Exit 0 (jeder Block gilt, Text außerhalb unbeachtet); Quelle Umgebung (`PGWIRE_RECORDER_UPSTREAM_CA` und `…_TLS`) und Quelle Schlüssel (`play: upstream_ca: ca.pem`, relativer Pfad ab dem aktuellen Verzeichnis): Exit 0 | bestätigt |
| Name gegen den eingesetzten Host | anderer Hostname (`pgssl2`, ein Server mit demselben Zertifikat): `PGR-E4005`, „certificate is valid for pgssl, localhost, not pgssl2“; andere IP (`172.31.77.11`): `PGR-E4005`, „valid for 172.31.77.10, not 172.31.77.11“. Beide Male fehlt die Tabelle | bestätigt |
| Abgelaufenes Zertifikat | `pgexp`: `PGR-E4005`, „certificate has expired or is not yet valid … after 2020-01-02“; Tabelle fehlt | bestätigt |
| Kein TLS ohne Wunsch, `sslmode=disable` und ohne `sslmode` | gegen `pgssl` (lehnt Klartext per `hostnossl reject` ab): `sslmode=disable` und ohne Angabe: `PGR-E4005`, „Fehlerantwort im Aufbau 28000 pg_hba.conf rejects connection … no encryption“, Exit 4. Das zeigt, dass `play` ohne den Wunsch im Klartext verbindet (der Server sah „no encryption“) | bestätigt |
| `--upstream-tls=false` gewinnt über `sslmode=require` aus jeder Quelle | Verbindung `plain` (`sslmode=require`) gegen `pgplain` (ohne SSL): `--upstream-tls=false`: Exit 0, Tabelle gefüllt; `PGWIRE_RECORDER_UPSTREAM_TLS=false`: Exit 0; Schlüssel `play: upstream_tls: false`: Exit 0. Gegenprobe: `--upstream-tls` (Option) gegen den Schlüssel `false`: `PGR-E4005` „lehnt TLS ab“ (die Option geht dem Schlüssel vor); Umgebung `true` gegen Option `false`: Exit 0 (die Kommandozeile geht der Umgebung vor); ohne alles: `PGR-E4005` | bestätigt |
| Kein Überspringen der Prüfung, keine weiteren `sslmode` | `--insecure`, `--upstream-insecure`, `--upstream-skip-verify`, `--sslmode`: je `PGR-E2001` („flag provided but not defined“); `sslmode=verify-full` und `prefer` in einer Verbindung: `PGR-E2004`, „erlaubt sind disable und require“. Der Quelltext enthält weder `InsecureSkipVerify` noch `MinVersion` (`grep`) | bestätigt |
| Integrationstest gegen einen Server mit dem Zertifikat einer eigenen CA | im eigenen `make gates` grün: `TestE2EPlayTLS` (5 Fälle), `TestE2EPlayTLSAnmeldung` (3), `TestE2EPlayTLSAbgelehnt` (7), `TestE2EPlayTLSZertifikatFehler` (5), `TestE2EPlayTLSZertifikatsspeicherDesSystems`, `TestE2EPlayTLSAbbruch`, `TestE2EPlayKlartextAbgelehnt`, `TestE2EPlayUpstreamCAFehler`. Das ist der Proxy der Tests; die Gegenstelle ist hier ein echter PostgreSQL 17 mit eigenem TLS (oben) | bestätigt |

### Punkt 2: Codes für `N`, Aushandlung, Zertifikat, Ablehnung ohne TLS (`PGR-E4005`); anderes Byte, Ende, Bytes nach `S` (`PGR-E4002`); `--upstream-ca` ohne TLS (`PGR-E2001`); CA-Datei (`PGR-E2007`, ohne Pfad und Inhalt); Beleg je Zusage Zusage · Mutation · roter Test. **Bestätigt, mit einer Einschränkung beim Beleg einer Zeile (V-167, LOW).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `N` auf das `SSLRequest`: `PGR-E4005`, **kein Klartext** | Fake-Server (`N`): `PGR-E4005` „lehnt TLS ab“; der Server sah genau die acht Bytes `00 00 00 08 04 d2 16 2f`, sendete `N`, danach kein Byte, nur das Ende der Verbindung („NACH N: 0 B“). Echter Server ohne SSL (`pgplain`) mit `sslmode=require` und mit `--upstream-tls`: `PGR-E4005`; sein Log zeigt für die beiden Läufe `connection received` und **kein** `connection authorized` (die nächsten `authorized`-Zeilen sind meine `psql`-Prüfungen als `postgres`). Es gab keine Klartext-Anmeldung | bestätigt |
| Anderes Byte: `PGR-E4002` | Fake-Server (`E`): `PGR-E4002`, „Antwort 'E' auf das SSLRequest ist weder S noch N“; danach kein Byte | bestätigt |
| Ende vor der Antwort: `PGR-E4002` | Fake-Server (schließt): `PGR-E4002`, „Verbindung vor der Antwort … beendet: EOF“ | bestätigt |
| Bytes nach `S`: `PGR-E4002` | Fake-Server (`Sxyz` in einem Segment): `PGR-E4002`, „Bytes des Servers nach S vor der Aushandlung“; danach kein Byte | bestätigt |
| Fehler der Aushandlung und des Zertifikats: `PGR-E4005` | siehe Punkt 1 (unbekannte CA, Name, IP, abgelaufen); Fake-Server (`S`, dann Schweigen): kein Fehler, `play` wartet (Signale unten) | bestätigt |
| Server lehnt Klartext ab: `PGR-E4005` | siehe Punkt 1 (28000, „no encryption“) | bestätigt |
| `--upstream-ca` ohne TLS: `PGR-E2001`, nach `--upstream` | `host:port` mit `--upstream-ca`: `PGR-E2001`, Meldung nennt die Option, nicht den Pfad; Verbindung mit `sslmode=require`, aber `--upstream-tls=false`: ebenso. Mit einer ungültigen `--input` zugleich: Datei zuerst (E2007 vor E3001, siehe unten) | bestätigt |
| CA-Datei: fehlend, Verzeichnis, leer, kein `CERTIFICATE`, nicht lesbar: `PGR-E2007`, ohne Pfad und Inhalt | je Fall Exit 2, `PGR-E2007`, Meldungen: „Datei nicht vorhanden“, „keine reguläre Datei“ (Verzeichnis und FIFO, die FIFO endet sofort), „kein PEM-Block“ (leer und nur Text), „ein PEM-Block hat nicht den Typ CERTIFICATE“ (Block `PRIVATE KEY`), „ein PEM-Block enthält kein lesbares X.509-Zertifikat“, „Datei nicht lesbar: permission denied“ (chmod 000). Die Dateien trugen `GEHEIMINHALT` und Base64 davon: `grep` nach Inhalt und nach den Pfaden (`/work/fehlt`, `/work/leer`, …) über alle Ausgaben: **0 Treffer**. Ein Link auf `ca.pem` wird gefolgt (Exit 0) | bestätigt |
| E2007 vor dem Laden der Aufzeichnung | `--input` fehlt, `--upstream-ca` leer: `PGR-E2007`, kein `PGR-E3001`; mit gültiger CA und fehlender Aufzeichnung: `PGR-E3001` | bestätigt |
| `config show` liest die Datei nicht (j) | `config show` mit `upstream_ca` auf eine unlesbare, eine fehlende, ein Verzeichnis und eine leere Datei (Schlüssel und Umgebungsvariable): Exit 0, zeigt Pfad der Konfiguration, den Inhalt unaufgelöst und den Namen `PGWIRE_RECORDER_UPSTREAM_CA` ohne Wert | bestätigt |
| `record` mit `sslmode=require` bleibt `PGR-E2004` (k) | Binary: `PGR-E2004` „sslmode=require ist bei record ungültig, record verbindet ohne TLS zum Upstream“, Exit 2; `record --upstream-tls` und `--upstream-ca`: `PGR-E2001` (unbekannte Option) | bestätigt |
| Anmeldung über TLS (h) | `play_scram`, `play_md5`, `play_pw` (Verfahren `scram-sha-256`, `md5`, `password`, nur über `hostssl`): je Exit 0, die Tabelle trägt `tls-ok`; der Server-Log zeigt `SSL enabled` für alle drei; falsches SCRAM-Passwort: `PGR-E4005` „Fehlerantwort im Aufbau 28P01“. Das Klartext-Passwort ging nur über TLS (der Server lehnt `host`-Zeilen ohne SSL ab). Kein `GEHEIM` in der Ausgabe von `--log-level debug` | bestätigt |
| Signale in der Aushandlung (i) | Fake-Server hält nach `S` (nach dem ClientHello) und Fake-Server antwortet nicht auf das `SSLRequest`: erstes `SIGINT`: `play` läuft weiter (das erste Signal lässt den Aufbau zu Ende laufen); zweites: Ende **0,22 s** und **0,19 s** später, Exit 0, Zeile „Abbruchsignal, play endet vorzeitig“, kein Fehler; der Server sah das Ende der Verbindung | bestätigt |
| Abbruch im Aufbau: kein `Terminate`, kein Klartext | Fake-Server (`N`, anderes Byte, Ende, `Sxyz`): jeweils nur das `SSLRequest`, danach nichts | bestätigt |
| Test je Zusage mit Mutation, roter Test (§3.10) | **37 Mutanten-Läufe** (36 im Unit-Image, 1 in der Integration; Tabelle unten): 34 rot aus dem richtigen Grund, ein Mutant bleibt wie in §7 angekündigt grün (V15, Redundanz im Produkt), zwei bleiben grün: V05 im Unit-Image (E2E rot, wie der Plan es verlangt) und **V32** (Behauptung in §7 nicht reproduzierbar, V-167) | bestätigt, V-167 für eine Zeile |

**Tabelle der Mutanten** (alle in frischer Kopie, gofmt-sauber; „rot“ heißt: ein genannter Test schlug mit seiner Meldung fehl, nicht Build oder Vet; drei erste Versuche waren ungültig und sind durch `b`-Fassungen ersetzt: V09 war äquivalent, V19 übersetzte nicht, bei V20 traf der Suchtext nicht).

| Gruppe | Mutant (Ersetzung) | roter Test |
|---|---|---|
| Prüfung des Zertifikats | V01 `InsecureSkipVerify: true`; V03 Kette geprüft, Name nicht (`VerifyPeerCertificate`); V25 Name fest `127.0.0.1`; V04 `pool.AddCert` entfernt; V14 Aushandlung als `PGR-E4002` | `TestEinspielTLSZertifikatFehler`, `TestEinspielTLSName`, `TestEinspielTLSZertifikatsspeicher`, `TestRunPlayTLS` (V01); `TestEinspielTLSZertifikatFehler`, `TestEinspielTLSName` (V03); `TestEinspielTLSName`, `TestEinspielTLSServerName`, `TestEinspielTLSIPv6MitZone` (V25); sieben Tests (V04); fünf Tests (V14) |
| Antwort auf `SSLRequest` | V02 `N` setzt im Klartext fort (Downgrade); V13 `N` als `PGR-E4002`; V12 anderes Byte als `PGR-E4005`; V29 Ende als `PGR-E4005`; V28 Senden als `PGR-E4005`; V06 `n > 2` statt `n > 1`; V31 nur ein Byte gelesen | `TestEinspielTLSAntwort` (V02: `/N` und `/N_mit_Bytes_dahinter`, Meldung „erwartet PGR-E4005“, Fehler bleibt `PGR-E4002` aus dem Klartext-Startup); `TestEinspielTLSAntwort`; `TestEinspielTLSAntwort`, `TestEinspielTLSSendenUndLesen`; `TestEinspielTLSAntwort`, `…SendenUndLesen`; `TestEinspielTLSSendenUndLesen`; `TestEinspielTLSAntwort`, `…SendenUndLesen`; dieselben zwei |
| Aufbau | V10 Startup auf der Klartext-Verbindung nach der Aushandlung; V11 TLS-Zweig aus; V21 Bootstrap reicht TLS nicht durch; V22 Bootstrap reicht CA nicht durch | `TestEinspielTLS`, `…KlartextPasswort`, `…AbbruchImAufbau` und drei weitere, `TestRunPlayTLS` (V10); neun Tests (V11); `TestRunPlayTLS` (V21, V22) |
| Vorrang und Quellen | V08 `sslmode` entscheidet immer; V09b `require` schlägt ein gesetztes `false`; V16 `--upstream-ca` ohne TLS erlaubt; V34 Quelle von `--upstream-ca` nicht gelesen | zehn Tests (V08); `TestPlayUpstreamTLS`, `TestPlayUpstreamCAOhneTLS` (V09b); `TestPlayUpstreamCAOhneTLS`, `…Reihenfolge` (V16, V34) |
| CA-Datei | V17 Typ des Blocks ungeprüft; V36 Typ `TRUSTED CERTIFICATE` zugelassen; V18 kein Block erlaubt; V30 nur der erste Block; V19b reguläre Datei ungeprüft (FIFO endet nicht binnen 10 s); V27 `Lstat` statt `Stat`; V20b Pfad in der Meldung; V23 Code `PGR-E2004` statt `PGR-E2007` | `TestPlayUpstreamCADateiFehler` (V17, V36, V23, V20b); `…Reihenfolge`, `…DateiFehler` (V18); `TestPlayUpstreamCADatei`, `…DateiFehler`, `…Quellen` (V30); `TestPlayUpstreamCAKeineRegulaere`, `…DateiFehler` (V19b); `TestPlayUpstreamCAKeineRegulaere` (V27) |
| Geheimnis | V24 `Format` von `cli.Zertifikate` gibt das Subjekt aus; V33 dasselbe im Upstream-Adapter | `TestPlayOptionenOhneZertifikat` (V24); `TestEinspielzielOhneZertifikat` (V33) |
| Systemspeicher | V26 nil-Speicher ohne Fehler; V05 Standardpfad leerer Speicher (M14) | `TestEinspielTLSZertifikatsspeicher` (V26); V05: Unit **grün**, **E2E** `TestE2EPlayTLSZertifikatsspeicherDesSystems/SSL_CERT_FILE_nennt_die_Zertifizierungsstelle` rot (Exit-Code und Meldung „unknown authority“) |
| Binary und Testhilfen | V07 Blank-Import von `tlsproxy` in `bootstrap.go`; V35 `close(p.Hello)` ohne `sync.Once` | `TestBinaryOhneTesthilfen` (V07); V35: Absturz des Testprozesses („close of closed channel“), `TestHaengenderProxyZweiteVerbindung` |
| Grün | V15 `HandshakeContext(context.Background())` (§7: `AfterFunc` schließt ebenfalls); V32 Fehler des Ladens nimmt `x509.SystemCertPool()` des Containers | beide ohne roten Test; V15 ist in §7 als Redundanz benannt, V32 siehe V-167 |

### Punkt 3: `make gates` grün. **Bestätigt durch eigenen Lauf.**

Eigener Lauf im Arbeitsbaum am Stand `a9c9330`, sauber vor und nach dem Lauf: Exit 0. Darin `baseline-verify: v6.18.0 OK — 54 Dateien`, `d-check: 520 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen — alle zwoelf Faelle`, `abdeckung-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `run-integration-tests: gruen` mit `--- PASS` für alle `TestE2EPlayTLS*`; kein `FAIL` im Log. Zusätzlich (Hinweis, kein Gate): `go test -race -count=8` über Upstream-Adapter, CLI-Adapter und `./internal/bootstrap/...` (Go 1.27.2 mit `gcc`): alle vier `ok`, keine Meldung.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`b221a16`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, Lerneintrag:** noch offen; §7 trägt „Wird bei Closure gefüllt“, beide Risiken aus §6 „offen bis Closure“. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Risiko 2 (Zertifikatsspeicher des Produkt-Images leer) ist nicht eingetreten: Das Image der Stufe `runtime` trägt `/etc/ssl/certs/ca-certificates.crt` mit 150 Zertifikaten, und `SSL_CERT_FILE` wirkt im Binary dieser Stufe. Dass ein Serverzertifikat einer öffentlichen CA im Betrieb geprüft wird, habe ich nicht gefahren (kein Netz, kein Zertifikat einer öffentlichen CA); das Risiko bleibt insoweit ein Urteil. Risiko 1 (Wachstum des Testgeschirrs: `testpki` und `tlsproxy` rund 250 Zeilen) ist Urteil des Planners. Aus der Verifikation gehören in den Lerneintrag die Klassen aus V-167 und V-168.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-601 Sicht nennt `go list` | `grep -rn 'go list' spec/` ohne Treffer; `spec/architecture.md` §6 sagt nur „darf nur Testcode importieren; das ausgelieferte Programm enthält sie nicht“ (kein Werkzeugaufruf, keine Wellen-, Slice- oder ADR-Kennung) | behoben; V-168 zum ersten Halbsatz |
| F-602 Zusage ohne Prüfung | `TestBinaryOhneTesthilfen` ruft `go list -deps ./cmd/...` ohne Netz auf und hält die Menge gegen `testing`, `testpki`, `tlsproxy`; Mutant V07 (Blank-Import von `tlsproxy`): **rot** mit „das Binary linkt …“. Beide `doc.go` sagen nur noch zu, was der Test prüft („Kein Paket unter ./cmd/... importiert es …; für Testdateien außerhalb … prüft das nichts“). Die Mutanten für `testing` und `testpki` allein habe ich nicht gefahren | behoben; V-168 |
| F-603 Standardpfad des Speichers | `TestE2EPlayTLSZertifikatsspeicherDesSystems` startet das Binary mit `SSL_CERT_FILE` im Prozessumfeld; die Gegenprobe ohne die Variable ist `PGR-E4005`. Mutant V05 (`SystemCertPool` durch leeren Speicher): E2E **rot** (Teilfall `SSL_CERT_FILE_nennt_die_Zertifizierungsstelle`), Unit grün, wie der Plan es sagt. Dasselbe am Binary der Stufe `runtime` mit dem echten Server (Punkt 1) | behoben |
| F-604 Name sagt mehr als der Test prüft | Der Fall heißt jetzt „S, dann Gegenstelle nur TLS 1.0, Aushandlung scheitert“ (`einspielen_tls_test.go` Zeile 438); `grep` nach `MinVersion` und `InsecureSkipVerify` im Quelltext: 0 Treffer; §6 führt die Version als Voreinstellung ohne Zusage, `LH-FA-20.a` *TLS* ebenso | behoben |
| F-605 Randformen im Code vor §6 | §6 trägt die vier Punkte als *Rückgabe des Implementers (geprüft vom Architect … nach dem Code)*; die Reihenfolge ist nicht mehr zu ändern | zur Kenntnis (Klasse für den Zähler) |
| F-606 `close(p.Hello)` | `sync.Once` in `proxy.go`; `TestHaengenderProxyZweiteVerbindung` öffnet zwei Verbindungen; Mutant V35 (ohne `Once`): **Absturz** des Testprozesses | behoben |
| F-607 Grenze am falschen Satzteil | `LH-FA-20.a` *TLS* sagt jetzt „Erkannt werden die Bytes nach S, die mit dem S eintreffen; was später eintrifft, stört die Aushandlung (Grenze: ohne Test …)“; Änderungsvermerk in der Tabelle der Spezifikation | behoben |
| F-608 Größe | der Plan benennt das in §7 (*Größe*); das Maß gehört dem Planner | zur Kenntnis |

---

## 3. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`; `a-check-negativ: gruen` (zwölf Fälle). [ADR-0001](../plan/adr/0001-hexagonale-architektur.md): `git diff --stat 3d0bd0c..a9c9330 -- .a-check.yml .golangci.yml Makefile harness cmd go.mod go.sum` ist leer; der Core (`internal/hexagon`) ändert sich nur um die Konstante `CodeConfigCA` (eine Zeile, `PGR-E2007`), kein `crypto/x509`, kein `crypto/tls`.
- [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md): `pgproto3` nur im Upstream-Adapter (`SSLRequest` über `pgproto3`, acht Bytes am Fake-Server gesehen: `00 00 00 08 04 d2 16 2f`); `crypto/tls` im Upstream-Adapter und in den Testhilfen unter `internal/bootstrap`, wo `.a-check.yml` es zulässt.
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md): Klartext, MD5 und SCRAM über TLS am echten Server gesehen; Zertifikat gegen den Systemspeicher; kein Überspringen (Punkt 1); keine Option für das Passwort. [ADR-0019](../plan/adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md): die Zertifikate aus `--upstream-ca` **ergänzen** den Systemspeicher (Datei mit zwei CAs, `SSL_CERT_FILE` ohne die Option; Mutant V04); `PGR-E2001` ohne TLS, `PGR-E2007` für die Datei. [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md): jeder Fehler im Aufbau verlässt den Adapter mit seinem Code (alle Läufe oben: `PGR-E4005` oder `PGR-E4002`, Exit 4, nie ein Fehler ohne Code).
- `spec/architecture.md` §4.4 (Authentifizierung und TLS des Upstream-Adapters, die Datei der CA liest der CLI-Adapter) und §6: Code und Sicht stimmen; der Satz zu den Testhilfen in §6 trägt keine Kennung und keinen Werkzeugaufruf.
- Kein `//nolint` im Diff (`make lint` im Gate grün, `lint-gegenprobe` grün); Profil und Regeln unverändert.
- Abdeckungstabellen: `make abdeckung-check` im Gate grün.

---

## 4. Plan gegen Code

- **§1:** Geliefert sind `--upstream-tls`, `--upstream-ca`, TLS nach `sslmode`, die Datei beim Start (`PGR-E2007`) und `PGR-E4005`/`PGR-E4002` im Aufbau. Die Abgrenzungen halten: Benutzerhandbuch und `README.md` sind nicht im Diff (Doku-Folge-Slice `slice-v1-abschluss-einspielen-tls-doku`); `record` ist unberührt (`PGR-E2004` bei `sslmode=require`, `record --upstream-tls` unbekannt); keine Client-Zertifikate, kein weiterer `sslmode`, kein Überspringen; kein Code im Core und im PGWire-Adapter.
- **§2 Zählung:** zwei Liefer-Punkte (DoD 1 und 2), zwei Schichten (Upstream-Adapter, CLI-Adapter); der Bootstrap setzt zwei Felder, `model` trägt eine Konstante. Das stimmt mit dem Diff (Bootstrap `+6 −1`).
- **§3:** Jede geänderte Datei steht dort, auch die Testhilfen und die Aufträge aus der Review. Die Zeile zu `internal/adapters/driving/cli` nennt `PlayOptions` „nicht mehr mit `==` vergleichbar“; `leser_test.go` vergleicht mit `reflect.DeepEqual`.
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort: *Optionen*, *`--upstream-tls` ausdrücklich `false`* (aus jeder Quelle), *`--upstream-ca` ohne TLS* (Reihenfolge), *Datei aus `--upstream-ca`* (reguläre Datei, Link gefolgt, FIFO ohne Warten, jeder Block `CERTIFICATE`, Text außerhalb unbeachtet, `PGR-E2007` ohne Pfad und Inhalt), *Antwort auf `SSLRequest`*, *Name im Zertifikat* (DNS und IP), *Host als IP-Adresse*, *Abbruchsignal in der Aushandlung*, *Server lehnt Klartext ab*, *Ablösung des Zwischenstands*. Nicht gefahren: *IPv6 mit Zone* (kein IPv6 im internen Netz; der Unit-Test `TestEinspielTLSIPv6MitZone` lief ohne Skip und ist mit V25 rot), *Speicher nicht ladbar* (nur über den Test-Hook beobachtbar, V-167), *Ablauf eines Zertifikats erst beim Aufbau* (gefahren: abgelaufen ist `PGR-E4005` beim Aufbau).
- **§7:** Die Zahl der Mutanten und die Tabelle stimmen im Wesentlichen. Der **Namensabgleich** (V-162-Lehre): Aus §7 zog ich alle Testfunktionen (43) und alle Teilfälle (rund 170 Namen in Backticks, Unterstrich als Leerzeichen) und suchte sie in allen Go-Dateien: alle gefunden bis auf `TestE2EPlayZwischenstand`, den §7 selbst als entfernt nennt; zwei Namen mit `SSL_CERT_FILE` stehen mit Leerzeichen im Quelltext (`t.Run("SSL_CERT_FILE nennt die Zertifizierungsstelle", …)`). Alle in §7 genannten Pfade und Commits existieren. Falsch ist eine Zeile der Grenzen (V-167).
- **§8:** Zähler und Sub-Area stimmen mit der Deklaration in `harness/conventions.md` (Sub-Area `REPO`, Greenfield).

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-167 | LOW | §7 behauptet in der Mutationstabelle, der Mutant „Fehler des Ladens nimmt die Systemwurzeln“ werde von `TestEinspielTLSZertifikatsspeicher/nicht_ladbarer_Speicher,_in_CA` rot, und in *Grenzen und nicht Geprüftes*, die Zeilen *leerer* und *nicht ladbarer Speicher* würden gegen diesen Mutanten „erst mit CA“ getrennt. In meiner Fassung des Mutanten (`pool, _ = x509.SystemCertPool()` im Fehlerzweig, danach CA hinzugefügt) bleibt das Paket grün (Mutant V32; ebenso die Pakete CLI und Bootstrap): Die Systemwurzeln des Test-Containers enthalten die Test-CA nicht, und mit CA gelingt der Aufbau in beiden Fassungen. Der Unterschied *leer* gegen *Systemwurzeln* ist nur mit einem Serverzertifikat einer vertrauten Wurzel zu beobachten, das die Testumgebung nicht erzeugen kann. Die Zusage steht in `LH-FA-20.a` *TLS* als „(Grenze)“; der Hook prüft nur *nil und Fehler gelten als leer* (V26 rot). | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls.md` · „Fehler des Ladens nimmt die Systemwurzeln“ und „getrennt wird das erst mit CA“ | Implementer: den Mutanten, der rot wird, in §7 genau benennen (Ersetzung), oder die Zeile und den Satz berichtigen („kein Test trennt leeren Speicher von den Systemwurzeln des Containers“) und die Grenze so führen; kein Code |
| V-168 | INFO | Die Sicht (`spec/architecture.md` §6) sagt „darf nur Testcode importieren“; `TestBinaryOhneTesthilfen` prüft den zweiten Halbsatz („das ausgelieferte Programm enthält sie nicht“), nicht den ersten: Ein Import durch eine Nicht-Testdatei außerhalb von `./cmd/...` prüft nichts (`doc.go` und §7 nennen das ausdrücklich als Grenze, `make a-check` nimmt `internal/bootstrap/**` aus). Bewusst so entschieden (Gate-Eingriff, `AGENTS.md` §3.6); die Sicht trägt die Norm, kein Sensor die Einhaltung. | `spec/architecture.md` · „darf nur Testcode importieren“; `internal/bootstrap/testhilfen_import_test.go` · „TestBinaryOhneTesthilfen“ | Architect: Kenntnis; ob die Norm der Sicht ohne Sensor stehen bleibt oder ein Gate (ADR) kommt, ist seine Entscheidung |

Die Falsifikationsversuche, die nichts fanden: Downgrade (`N` und andere Antworten, Bytes nach `S`): am Server in keinem Fall ein Byte nach dem `SSLRequest`, im Log von `pgplain` zweimal `connection received` ohne `connection authorized`; die Reihenfolge der Quellen von `--upstream-tls` (Option, Umgebung, Schlüssel gegen `sslmode`) in allen Kombinationen; `PGR-E2007` ohne Pfad und Inhalt in neun Fällen; `PGR-E2001` und `PGR-E2007` vor dem Laden der Aufzeichnung; Mehrfach-CA und Text außerhalb der Blöcke; Passwort in keiner Ausgabe (`--log-level debug`, Fehlerfälle); keine Option zum Überspringen; Signale in der Aushandlung; `-race -count=8` ohne Meldung; ein Name- und IP-Abgleich gegen einen unabhängigen TLS-Server, den nicht der eigene Code gebaut hat.

Befunde: 0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO. Alle drei DoD-Verhaltenspunkte sind bestätigt, auch gegen eine Gegenstelle, die nicht der eigene Code ist; der Beleg einer Zeile in §7 (V-167) ist nicht reproduzierbar, ohne dass eine Zusage des Produkts betroffen wäre.

**Finding-Klassen dieses Laufs:** Mutant im Beleg benannt, in meiner Ersetzung grün · Norm der Sicht ohne Sensor für den ersten Halbsatz

**Closure-blockierend:** nein. V-167 ist ein Nachzug an §7 (eine Zeile und ein Satz), keine Änderung am Code.

**Übergabe an den Planner:** dieser Bericht mit V-167 (an den Implementer) und V-168 (an den Architect); dazu für die Closure die Ausgänge der zwei Risiken aus Abschnitt 1 (*Übrige DoD-Punkte*) und die Klassen oben. Für die Doku-Folge: Das Handbuch nennt `sslmode=require` und `--upstream-tls` bei `play` laut §7 weiter als nicht verfügbar (`slice-v1-abschluss-einspielen-tls-doku`).

**Aufgeräumt:** Container `ver-tls-pgssl`, `-pgssl2`, `-pgexp`, `-pgplain`, `-fake`, `-rec`, `-sig` und die Läufe `ver-tls-it-*`; Netz `ver-tls-net`; Volumes `ver-tls-gocache` und `ver-tls-gomod`; Images `ver-tls-rt`, `ver-tls-deps`, `ver-tls-int`; die Verzeichnisse `src`, `mut`, `work`, `certs` im Scratchpad. `docker ps -a`, `docker network ls`, `docker volume ls` und `docker images` mit `ver-tls` danach leer; die Images des Arbeitsbaums (`pgwire-recorder:*`) blieben stehen.
