# Benutzerhandbuch: pgwire-recorder

Version: 0.2  
Software-Version: `dev` (Ausgabe von `pgwire-recorder version`)  
Stand: 11.10.2026  
Autor: Projektteam pgwire-recorder  
Gültigkeitsbereich: gilt für das Binary `pgwire-recorder`, gebaut aus dem Repository (Ausgabe von `pgwire-recorder version`: `pgwire-recorder dev`)

## 1. Einleitung

### Zweck der Software

`pgwire-recorder` zeichnet die Kommunikation zwischen Ihrer Anwendung und einer
PostgreSQL-Datenbank auf und spielt sie später ohne die Datenbank wieder ab.
Damit laufen Tests Ihrer Anwendung reproduzierbar, ohne dass dafür eine
PostgreSQL-Instanz bereitstehen muss.

Das Werkzeug arbeitet in drei Betriebsarten:

* **Aufzeichnen (`record`):** Ihre Anwendung verbindet sich mit
  `pgwire-recorder` statt direkt mit PostgreSQL. Das Werkzeug leitet alles an
  die echte Datenbank weiter und speichert die Kommunikation in einer Datei.
* **Wiedergeben (`replay`):** Ihre Anwendung verbindet sich wieder mit
  `pgwire-recorder`. Das Werkzeug beantwortet die Anfragen aus der Datei. Eine
  Datenbank ist nicht nötig.
* **Einspielen (`play`):** Das Werkzeug führt die aufgezeichneten Anfragen aus der
  Datei gegen eine Datenbank aus. Ihre Anwendung ist nicht beteiligt.

### Zielgruppe

Dieses Handbuch richtet sich an Entwicklerinnen, Entwickler und Tester, die
Anwendungen mit PostgreSQL-Anbindung testen, sowie an Teams, die diese Tests in
CI-Systemen ausführen.

### Voraussetzungen

* Ihre Anwendung kann Host und Port der Datenbankverbindung einstellen. Weitere
  Änderungen an der Anwendung sind nicht nötig.
* Zum Aufzeichnen: eine erreichbare PostgreSQL-Datenbank, die den Benutzer ohne
  Passwort anmeldet. Verlangt sie ein Passwort, beendet `record` die Verbindung
  Ihrer Anwendung mit `PGR-E6001`.
* Zum Einspielen: eine erreichbare PostgreSQL-Datenbank. Verlangt sie ein Passwort,
  nennen Sie es `play` (siehe [Einstellungen](#5-einstellungen)). Auf Wunsch verbindet
  `play` mit TLS (siehe
  [Eine Aufzeichnung in eine Datenbank einspielen](#eine-aufzeichnung-in-eine-datenbank-einspielen)).
* Zum Bauen: Docker und GNU `make`. Zum Ausführen: Linux auf der Architektur des
  Rechners, auf dem Sie gebaut haben, oder ein Container-Laufzeitsystem (siehe
  [Installation](#2-installation)).
* Die Verbindung zum Werkzeug läuft ohne Verschlüsselung: `record` und `replay`
  lehnen eine Anfrage nach Verschlüsselung ab. Ein Treiber mit `sslmode=prefer`
  verbindet sich dann unverschlüsselt, einer mit `sslmode=require` bricht ab (siehe
  [Fehlerbehebung](#die-anwendung-kann-sich-nicht-verbinden)).

### Wichtig: Aufzeichnungen können vertrauliche Daten enthalten

Eine Aufzeichnung enthält die SQL-Anfragen und die zurückgegebenen Datensätze.
Das Werkzeug erkennt oder maskiert vertrauliche Inhalte nicht. Verwenden Sie
zum Aufzeichnen geeignete Testdaten, und speichern, übertragen und versionieren
Sie Aufzeichnungen so, wie es zu ihrem Inhalt passt.

## 2. Installation

Das Werkzeug entsteht aus dem Repository des Projekts. Der Build läuft in Docker;
auf Ihrem Rechner brauchen Sie dafür nur Docker und GNU `make`.

### Aus dem Repository bauen

1. Klonen Sie das Repository, und wechseln Sie in sein Verzeichnis.
2. Führen Sie `make build` aus. Das Ergebnis ist das Image `pgwire-recorder:dev`;
   es enthält das Binary und braucht keine weitere Software.
3. Prüfen Sie den Build:

   ```bash
   docker run --rm pgwire-recorder:dev version
   ```

   Die Ausgabe ist `pgwire-recorder dev`.

### Das Binary aus dem Image

Das Binary im Image ist ein statisch gebundenes Linux-Programm für die Architektur
des Rechners, auf dem Sie gebaut haben. So kopieren Sie es in das aktuelle
Verzeichnis:

```bash
docker create --name pgwire-recorder-binary pgwire-recorder:dev
docker cp pgwire-recorder-binary:/pgwire-recorder ./pgwire-recorder
docker rm pgwire-recorder-binary
./pgwire-recorder version
```

Die Beispiele dieses Handbuchs rufen das Binary als `pgwire-recorder` auf. Legen Sie
es dafür in ein Verzeichnis Ihres Suchpfads.

### Container

Im Image ist das Binary der Einstiegspunkt: Die Argumente nach dem Image-Namen sind
Kommando und Optionen. Das Image läuft als Benutzer `nonroot` (UID 65532); ein
eingehängtes Verzeichnis, in das `record` schreibt, muss für ihn beschreibbar sein.
Ein Beispiel für eine Wiedergabe in Docker Compose, mit dem Image aus
`make build` auf demselben Rechner:

```yaml
services:
  recorder:
    image: pgwire-recorder:dev
    command:
      - replay
      - --listen=0.0.0.0:5432
      - --input=/recordings/test.yaml
    volumes:
      - ./recordings:/recordings:ro
```

## 3. Erste Schritte

Das folgende Beispiel zeichnet eine Anwendung auf und spielt sie danach ohne
Datenbank wieder ab.

1. Starten Sie das Aufzeichnen:

   ```bash
   pgwire-recorder record \
     --listen 0.0.0.0:15432 \
     --upstream postgres:5432 \
     --output ./recordings/users.yaml
   ```

2. Richten Sie Ihre Anwendung auf `localhost:15432` aus, statt auf die
   Datenbank, und führen Sie den Ablauf aus, den Sie aufzeichnen möchten.
3. Beenden Sie das Werkzeug mit `Strg+C`. Dabei wird die Datei
   `./recordings/users.yaml` geschrieben.
4. Beenden Sie die Datenbank, oder lassen Sie sie ungenutzt.
5. Starten Sie das Wiedergeben:

   ```bash
   pgwire-recorder replay \
     --listen 0.0.0.0:15432 \
     --input ./recordings/users.yaml
   ```

6. Führen Sie denselben Ablauf mit Ihrer Anwendung erneut aus, wieder gegen
   `localhost:15432`.

Ihre Anwendung erhält dieselben Antworten wie beim Aufzeichnen.

## 4. Aufgaben ausführen

### Eine Anwendung aufzeichnen

#### Voraussetzung

Die Datenbank ist erreichbar, und die Zieldatei existiert noch nicht.

#### Vorgehen

1. Starten Sie `pgwire-recorder record` mit den Optionen `--listen`,
   `--upstream` und `--output`.
2. Verbinden Sie Ihre Anwendung mit der Adresse aus `--listen`.
3. Führen Sie den gewünschten Ablauf aus.
4. Beenden Sie das Werkzeug mit `Strg+C` oder `SIGTERM`.

#### Ergebnis

Die Datei aus `--output` enthält die aufgezeichnete Kommunikation. Sie wird nach
dem Ende jeder Verbindung und beim Beenden aktualisiert.

#### Hinweise

* Existiert die Zieldatei bereits, bricht das Werkzeug ab (`PGR-E2002`). Wollen
  Sie sie ersetzen, ergänzen Sie `--force`; dasselbe gilt für
  `PGWIRE_RECORDER_FORCE=true` und für `force: true` im Abschnitt `record` der
  Konfigurationsdatei.
* Ist `--output` ein Verzeichnis oder sonst keine reguläre Datei, bricht das
  Werkzeug beim Start ab (`PGR-E3001`), auch mit `--force`. Das Verzeichnis der
  Zieldatei muss bestehen und beschreibbar sein; das Werkzeug legt es nicht an und
  bricht sonst beim Start ab (`PGR-E3001`), auch mit `--force`.
* Ist `--output` eine symbolische Verknüpfung, zählt ihr Ziel: Zeigt sie auf eine
  vorhandene Datei, gilt die Zieldatei als vorhanden. Mit `--force` ersetzt das
  Werkzeug die Verknüpfung durch die Aufzeichnung; ihr Ziel bleibt unverändert.
* Das Werkzeug schreibt die Aufzeichnung zuerst in
  eine temporäre Datei `.<Name der Zieldatei>.<16 Hexziffern>.tmp` im Verzeichnis
  der Zieldatei und ersetzt die Zieldatei danach in einem Schritt. Die Zieldatei
  ist damit vollständig oder unverändert, auch nach dem zwangsweisen Beenden von
  Verbindungen beim Herunterfahren. Schlägt das Schreiben fehl, meldet das
  Werkzeug `PGR-E3001` und entfernt die temporäre Datei; gelingt das Entfernen
  nicht, nennt dieselbe Meldung auch diesen Fehler. Endet der Prozess hart
  (zum Beispiel mit `SIGKILL`), kann sie liegen bleiben; das Werkzeug entfernt
  eine solche Datei eines früheren Laufs nicht und meldet sie nicht, entfernen
  Sie sie selbst.
* Eine neue Zieldatei erhält die Rechte `0666` nach der umask des Prozesses (bei
  der umask `022` also `0644`); eine ersetzte behält ihre Zugriffsrechte. Ersetzt
  das Werkzeug eine Verknüpfung auf eine Datei, erhält die Zieldatei die
  Zugriffsrechte dieser Datei.
* Jede Verbindung Ihrer Anwendung mit mindestens einer Anfrage wird als eigene
  Sitzung aufgezeichnet. Verbindungen ohne Anfrage, zum Beispiel
  Probe-Verbindungen eines Connection-Pools, werden nicht aufgezeichnet.
* Verlangt die Datenbank ein Passwort, beendet das Werkzeug die Verbindung Ihrer
  Anwendung mit `PGR-E6001`; die Verbindung wird nicht aufgezeichnet, und das
  Werkzeug endet mit Exit-Code 6.
* Beim Beenden nimmt das Werkzeug keine neuen Verbindungen an, schreibt eine
  Log-Zeile mit `sessions`, der Zahl der noch offenen Verbindungen, und wartet,
  bis jede Verbindung ihre laufende Anfrage oder Folge abgeschlossen hat. Es
  wartet höchstens so lange, wie `--shutdown-timeout` angibt (siehe
  [Herunterfahren mit Frist](#herunterfahren-mit-frist)); danach beendet es jede
  noch laufende Verbindung zwangsweise. Die abgeschlossenen Anfragen einer
  solchen Verbindung bleiben in der Aufzeichnung. War auf ihr eine Anfrage
  begonnen und nicht abgeschlossen, fehlt diese; Ihre Anwendung erhält die
  Fehlermeldung `PGR-E4006`, das Log nennt die Sitzung und die verworfene
  Anfrage, und das Werkzeug endet mit Exit-Code 4, wenn vorher kein anderer
  Fehler gemerkt wurde (siehe [Exit-Codes](#exit-codes)). Eine Verbindung ohne
  begonnene Anfrage, auch eine noch im Aufbau, endet ohne Meldung. Danach
  schreibt das Werkzeug die Aufzeichnung.
* Im Container beendet das Laufzeitsystem das Werkzeug nach seiner Stopp-Frist
  hart (bei Docker 10 Sekunden; `SIGKILL`, Exit-Code 137), dann fehlt die
  Aufzeichnung jeder noch wartenden Verbindung. Der Standardwert von
  `--shutdown-timeout` (5 Sekunden) liegt darunter, damit die Aufzeichnung
  geschrieben wird. Wählen Sie eine längere Frist nur mit längerer Stopp-Frist,
  in Docker Compose zum Beispiel mit `stop_grace_period`; das Verzeichnis
  `./recordings` muss für den Benutzer des Images beschreibbar sein (siehe
  [Container](#container)):

  ```yaml
  services:
    recorder:
      image: pgwire-recorder:dev
      command:
        - record
        - --listen=0.0.0.0:5432
        - --upstream=postgres:5432
        - --output=/recordings/test.yaml
        - --shutdown-timeout=30s
      volumes:
        - ./recordings:/recordings
      stop_grace_period: 60s
  ```

### Eine Aufzeichnung wiedergeben

#### Voraussetzung

Eine Aufzeichnung mit einer Verbindung liegt vor.

#### Vorgehen

1. Starten Sie `pgwire-recorder replay` mit den Optionen `--listen` und
   `--input`.
2. Verbinden Sie Ihre Anwendung mit der Adresse aus `--listen`.
3. Führen Sie denselben Ablauf aus wie beim Aufzeichnen.

#### Ergebnis

Ihre Anwendung erhält die aufgezeichneten Antworten, in der aufgezeichneten
Reihenfolge, auch Fehlerantworten der Datenbank.

#### Hinweise

* Die Wiedergabe ist streng: Jede Anfrage muss genau der aufgezeichneten
  Anfrage an dieser Stelle entsprechen, Zeichen für Zeichen (beim erweiterten
  Protokoll auch in Namen, Parametern, Formaten und Zeilenlimit). Eine abweichende
  Anfrage wird als Fehler gemeldet; das Werkzeug liefert dann keine geratene
  Antwort (`PGR-E5001`). Ausgenommen sind Anfragen, die nur aus Leerraum und
  Kommentaren bestehen, wie Treiber sie als Lebendprüfung senden (siehe
  [Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)).
* Die erste Verbindung mit einer Anfrage erhält die erste aufgezeichnete Sitzung,
  die zweite die zweite, und so weiter. Verbindungen ohne Anfrage zählen nicht.
  Nutzen Sie die Verbindungen nacheinander; bei gleichzeitiger Nutzung ist die
  Zuordnung nicht festgelegt. Eine Anfrage über die aufgezeichneten Sitzungen
  hinaus wird als Abweichung gemeldet (`PGR-E5003`).
* Endet eine Verbindung, bevor alle Anfragen ihrer Sitzung gestellt und
  beantwortet sind, warnt das Werkzeug (`PGR-W2001`), gleich ob die Anwendung
  die Verbindung schließt, ein Fehler sie beendet oder das Werkzeug beim
  Beenden. Eine Folge des erweiterten Protokolls zählt erst mit der Antwort auf
  ihr letztes `Sync` als gestellt. Die Meldung nennt die Sitzung, wie viele
  ihrer Anfragen offen sind, wie viele sie hat, und die Nummer der ersten
  offenen; Lebendprüfungen zählen nicht mit. Aufgezeichnete Sitzungen, die nie
  eine Verbindung erhalten haben, meldet das Werkzeug beim Beenden mit ihrer
  Zahl und der Kennung der ersten. Mit `--fail-on-unconsumed` ist beides ein
  Fehler (`PGR-E5002`, Exit-Code 5) mit demselben Text; die Anwendung erhält
  ihn nicht, weil die Verbindung dann schon endet.
* Beim Beenden beantwortet das Werkzeug eine begonnene Folge des erweiterten
  Protokolls noch bis zu ihrem `Sync` und wartet dabei höchstens so lange, wie
  `--shutdown-timeout` angibt (siehe
  [Herunterfahren mit Frist](#herunterfahren-mit-frist)). Pausiert Ihre
  Anwendung länger, schließt es die Verbindung. War eine Anfrage begonnen und
  nicht vollständig beantwortet, erhält die Anwendung die Fehlermeldung
  `PGR-E4006`, das Log nennt Sitzung und Anfrage, und das Werkzeug endet mit
  Exit-Code 4, wenn vorher kein anderer Fehler gemerkt wurde; mit
  `--fail-on-unconsumed` zählt `PGR-E4006` vor `PGR-E5002`.
* Im Container gilt die Stopp-Frist des Laufzeitsystems wie beim Aufzeichnen.
  Hinweis: Der Standardwert von `--shutdown-timeout` (5 Sekunden) ist so
  gewählt, dass er unter der Stopp-Frist von Docker (10 Sekunden) liegt; wählen
  Sie eine längere Frist nur mit längerer Stopp-Frist.

### Eine Aufzeichnung in eine Datenbank einspielen

Damit führen Sie die aufgezeichneten Anfragen erneut aus, zum Beispiel um eine
Komponente zu testen, die Änderungen der Datenbank verarbeitet (Change Data Capture).

#### Voraussetzung

Eine Aufzeichnung liegt vor, und die Zieldatenbank ist
erreichbar. Verlangt sie ein Passwort (Klartext, MD5 oder SCRAM-SHA-256), nennen Sie
es `play` über den Platzhalter der Verbindung oder die Umgebungsvariable
`PGWIRE_RECORDER_PASSWORD` (siehe [Einstellungen](#5-einstellungen)). Der
Benutzer und die Datenbank aus der Aufzeichnung existieren dort oder Sie geben sie
ausdrücklich an.

#### Vorgehen

1. Starten Sie das Einspielen:

   ```bash
   pgwire-recorder play \
     --upstream postgres:5432 \
     --input ./recordings/users.yaml
   ```

2. Warten Sie, bis das Werkzeug endet.

Mit TLS gegen einen Server, dessen Zertifikat eine eigene Zertifizierungsstelle
ausgestellt hat: Die Konfigurationsdatei `tls.yaml` wünscht TLS mit `sslmode=require`,

```yaml
connections:
  tls: "postgresql://dev@localhost:5432/myapp?sslmode=require"
play:
  upstream: tls
  input: ./recordings/users.yaml
```

und `--upstream-ca` nennt die Datei mit dem Zertifikat dieser Zertifizierungsstelle:

```bash
pgwire-recorder play --config ./tls.yaml --upstream-ca ./ca.pem
```

#### Ergebnis

Die Anfragen der Aufzeichnung sind in der aufgezeichneten Reihenfolge gegen die
Datenbank ausgeführt; sie enthält danach deren Wirkung. Jede Sitzung der
Aufzeichnung läuft über eine eigene Verbindung, die Sitzungen nacheinander.

#### Hinweise

* Antwortet die Datenbank auf eine Anfrage mit einem Fehler, bricht das Einspielen
  ab (`PGR-E4004`, Exit-Code 4). Mit `--continue-on-error` läuft es mit der
  nächsten Anfrage weiter und meldet jeden dieser Fehler im Log; der Exit-Code ist
  dann 4. Bricht
  danach ein anderer Fehler das Einspielen ab, etwa eine Antwort, die das Werkzeug
  nicht verarbeiten kann (`PGR-E6001`), oder das Ende der Verbindung
  (`PGR-E4003`), ist der Exit-Code der dieses Fehlers (6 bzw. 4).
* Antwortet die Datenbank auf eine Anfrage des erweiterten Protokolls mit einem
  Fehler, führt sie die übrigen Nachrichten dieser Folge bis zu deren `Sync` nicht
  aus. Mit `--continue-on-error` macht das Einspielen danach mit der nächsten
  Folge weiter, der Exit-Code ist dann 4.
* Mit `--allow-recorded-errors` gilt ein Fehler der Datenbank mit dem Schweregrad
  `ERROR` nicht als Fehler, wenn auch die aufgezeichnete Anfrage mit einem Fehler
  beantwortet wurde. Bei einer Folge des erweiterten Protokolls genügt eine
  Fehlerantwort an irgendeiner Stelle der aufgezeichneten Folge. Ein Fehler mit dem
  Schweregrad `FATAL` beendet die Verbindung
  und bricht das Einspielen auch dann ab (`PGR-E4003`).
* Ohne Wunsch nach TLS verbindet sich das Werkzeug ohne TLS mit der Datenbank. Lehnt
  sie die Anmeldung ab, etwa weil der Benutzer fehlt oder das Passwort falsch ist,
  oder fehlt das Passwort, das sie verlangt, endet das Einspielen mit `PGR-E4005`
  (Exit-Code 4); ebenso, wenn die Datenbank Verbindungen ohne Verschlüsselung ablehnt.
* TLS wünschen Sie mit `--upstream-tls` oder mit `sslmode=require` in der benutzten
  Verbindung (siehe [Konfigurationsdatei](#konfigurationsdatei)). Ein gesetztes
  `--upstream-tls` geht `sslmode` vor, auch `--upstream-tls=false`, gleich ob es als
  Option, als Umgebungsvariable oder als Schlüssel `upstream_tls` im Abschnitt `play:`
  kommt. Ohne beides verbindet `play` ohne TLS. Mit TLS baut jede Sitzung der
  Aufzeichnung ihre Verbindung mit TLS auf.
* Mit TLS prüft `play` das Zertifikat des Servers; ein Überspringen der Prüfung gibt es
  nicht. Das Zertifikat gilt, wenn eine Zertifizierungsstelle aus seiner Kette im
  Zertifikatsspeicher des Systems oder in der Datei von `--upstream-ca` steht (die Kette
  besteht aus dem Zertifikat des Servers und den Zwischenzertifikaten, die der Server
  mitsendet), und wenn es zum Host der Verbindung passt, so wie das Werkzeug ihn einsetzt: ein Hostname
  steht als DNS-Name, eine IPv4- oder IPv6-Adresse als IP-Adresse im Zertifikat. Es
  darf außerdem nicht abgelaufen sein. In diesen Fällen endet `play` mit `PGR-E4005` (Exit-Code
  4): Das Zertifikat des Servers lässt sich weder mit dem Speicher des Systems noch mit
  `--upstream-ca` überprüfen, der Host steht nicht im Zertifikat, das Zertifikat ist
  abgelaufen, oder der Server antwortet auf die Anfrage nach Verschlüsselung mit `N`.
  Auf `N` folgt kein Rückfall auf eine Verbindung ohne TLS. Die Meldung nennt die
  Sitzung und den Grund, nie Pfad oder Inhalt der Datei aus `--upstream-ca`.
  Antwortet der Server auf die Anfrage nach Verschlüsselung weder mit `S` noch mit `N`
  oder beendet er die Verbindung davor, endet `play` mit `PGR-E4002`.
* `--upstream-ca <Datei>` nennt eine PEM-Datei mit einem oder mehreren Zertifikaten
  von Zertifizierungsstellen. Sie ergänzt den Zertifikatsspeicher des Systems und
  ersetzt ihn nicht. Die Option verlangt, dass `play` TLS nutzt (sonst `PGR-E2001`): mit `--upstream-tls` oder `sslmode=require` der Verbindung, und ein gesetztes `--upstream-tls=false` heißt, dass `play` kein TLS nutzt. Ist die
  Datei nicht lesbar, keine reguläre Datei oder enthält sie kein Zertifikat, startet
  `play` nicht (`PGR-E2007`, Exit-Code 2); die Meldung nennt die Option und den Grund,
  nicht den Pfad. Die Zertifikate in der Datei, die die Umgebungsvariable
  `SSL_CERT_FILE` nennt, zählen zum Zertifikatsspeicher des Systems.
* Verlangt die Datenbank ein Klartext-Passwort, sendet `play` es auf der Verbindung,
  wie sie ist: mit TLS auf der verschlüsselten, ohne TLS auf der unverschlüsselten.
* Antwortet die Datenbank mit einem COPY-Datenstrom (`COPY … FROM STDIN`,
  `COPY … TO STDOUT`), kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001`
  (Exit-Code 6); die Meldung nennt Sitzung und Nummer der Anfrage.
* Bei `Strg+C` oder `SIGTERM` endet das Einspielen nach der laufenden Anfrage, bei
  einer Folge des erweiterten Protokolls nach deren `Sync`; mit
  `--finish-session-on-interrupt` erst nach der laufenden Sitzung. Ein zweites
  Signal schließt die Verbindung sofort, auch mit dieser Option; die unterbrochene
  Anfrage zählt nicht als Fehler. Der Exit-Code ist 0, nach einem früheren Fehler 4.
* Das Einspielen führt Anfragen aus. Verwenden Sie es nicht gegen eine Datenbank,
  deren Inhalt Sie nicht verändern dürfen.

### Die gewählte Konfigurationsdatei anzeigen

#### Voraussetzung

Eine Konfigurationsdatei liegt vor, zum Beispiel `.pgwire-recorder.yaml` im aktuellen
Verzeichnis.

#### Vorgehen

1. Führen Sie `pgwire-recorder config show` aus, bei Bedarf mit `--config <Datei>`.

#### Ergebnis

Das Werkzeug nennt in der ersten Zeile die gewählte Datei und gibt danach ihren
Inhalt aus, mit zwei Leerzeichen Einzug, in der Reihenfolge der Datei und ohne
Kommentare; Platzhalter `${VAR}` erscheinen unaufgelöst. Am Ende stehen die Namen
der gesetzten `PGWIRE_RECORDER_*`-Umgebungsvariablen, nach Namen sortiert und ohne
Wert, auch solche, zu denen es keine Option gibt. Der Befehl endet mit Exit-Code 0.
Findet er keine Datei, sagt er das in der ersten Zeile und endet ebenfalls mit
Exit-Code 0.

#### Hinweise

* Die Ausgabe enthält nie einen aufgelösten Wert, aber Hosts, Benutzer und Pfade der
  Datei. Prüfen Sie sie, bevor Sie sie weitergeben.
* Ist die Datei ungültig, zeigt das Werkzeug nichts und meldet `PGR-E2004` oder
  `PGR-E2006`. `PGR-E2005` meldet `config show` nie, weil es keine Verbindung benutzt.

### Mit einem Datenbanktreiber arbeiten

`record` und `replay` unterstützen das einfache und das erweiterte Anfrageprotokoll
von PostgreSQL. Viele Treiber nutzen standardmäßig das erweiterte Protokoll mit
vorbereiteten Anweisungen. Dafür müssen Sie am Treiber nur Host und Port ändern.

#### Voraussetzung

Ihr Treiber verwendet Version 3.0 des PostgreSQL-Protokolls und lässt sich ohne
Verschlüsselung betreiben.

#### Vorgehen

1. Starten Sie `pgwire-recorder record` oder `pgwire-recorder replay`.
2. Tragen Sie in der Verbindungszeichenfolge Ihres Treibers Host und Port aus
   `--listen` ein.
3. Stellen Sie die Verschlüsselung in der Verbindungszeichenfolge ab, zum Beispiel
   mit `sslmode=disable`; mit `sslmode=prefer` verbindet sich der Treiber
   unverschlüsselt.
4. Führen Sie Ihren Ablauf aus.

#### Ergebnis

Ihre Anwendung läuft wie gegen die Datenbank. Beim Aufzeichnen enthält die Datei
alle Anfragen, beim Wiedergeben erhält Ihre Anwendung die aufgezeichneten
Antworten.

#### Hinweise

* Beim erweiterten Protokoll müssen bei der Wiedergabe die Namen vorbereiteter
  Anweisungen und Portale, die SQL-Texte, die Parametertypen, die Parameterwerte,
  die Formate und das Zeilenlimit mit der Aufzeichnung übereinstimmen. Erzeugt
  ein Treiber bei jedem Lauf andere Namen, meldet das Werkzeug eine Abweichung
  (`PGR-E5001`).
* Weicht ein Parameterwert ab, nennt die Meldung die Nummer des ersten
  abweichenden Parameters, gezählt ab 1 wie `$1`, etwa `abweichend in params
  (Parameter $2)`; weicht die Zahl der Parameter ab, nennt sie beide Anzahlen.
  Die Werte selbst erscheinen nie, weder im Log noch in der Fehlermeldung an die
  Anwendung, auch nicht mit `--log-level debug`.
* Sendet ein Treiber nach einer Anforderung `Flush` weitere Nachrichten, ohne auf
  die Antwort zu warten, kann sich die Aufzeichnung zwischen zwei Läufen
  unterscheiden; die Wiedergabe einer vorhandenen Aufzeichnung bleibt davon
  unberührt.
* Beim Aufzeichnen beendet das Werkzeug die Verbindung mit `PGR-E6001`, wenn ein
  Treiber eine einfache Anfrage sendet, während eine Folge des erweiterten
  Protokolls noch auf ihr `Sync` wartet; die Verbindung wird dann nicht
  aufgezeichnet.
* Beim Beenden wartet das Werkzeug, bis eine laufende Folge des erweiterten
  Protokolls mit ihrem `Sync` abgeschlossen ist; eine neue Anfrage oder Folge
  beginnt danach nicht mehr, und die Verbindung wird geschlossen. Schließen Sie
  solche Folgen ab, bevor Sie das Werkzeug beenden.
* Connection-Pools und `database/sql` prüfen eine Verbindung, die eine Weile
  geruht hat, vor der nächsten Nutzung mit einer Lebendprüfung: einer Anfrage,
  die nur aus einem Kommentar besteht, bei pgx `-- ping` nach mehr als einer
  Sekunde Ruhe. Ob sie kommt, hängt also davon ab, wie schnell Ihr Ablauf
  läuft. Die Wiedergabe beantwortet eine solche Anfrage zwischen zwei Anfragen
  Ihres Ablaufs wie die Datenbank, ohne die Aufzeichnung zu verbrauchen, und
  überspringt aufgezeichnete Lebendprüfungen, die beim Wiedergeben ausbleiben.
  Sie müssen die Lebendprüfung am Treiber deshalb nicht abschalten (bei pgx
  `ShouldPing`). Das gilt für jede Anfrage, die nur aus Leerraum und
  Kommentaren besteht; eine Anfrage mit einer Anweisung, auch ein einzelnes `;`,
  muss weiter der Aufzeichnung entsprechen. Einen vertikalen Tabulator zählt die
  Wiedergabe wie PostgreSQL erst ab Version 17 zum Leerraum, nach der Version der
  aufgezeichneten Datenbank. Vor der ersten Anfrage einer Verbindung zählt er nie
  dazu: Eine Anfrage mit vertikalem Tabulator wird dort wie jede andere mit der
  Aufzeichnung verglichen und ist meist eine Abweichung (`PGR-E5001`, ohne freie
  Sitzung `PGR-E5003`). Mitten in einer Folge des
  erweiterten Protokolls, vor ihrem `Sync`, ist auch eine Lebendprüfung eine
  Abweichung (`PGR-E5001`). Aufzeichnen und Einspielen behandeln
  Lebendprüfungen wie jede andere Anfrage.

### Das Werkzeug in der Testautomatisierung verwenden

Das Werkzeug fragt nie nach Eingaben und lässt sich über Optionen und
Umgebungsvariablen vollständig steuern (siehe [Einstellungen](#5-einstellungen)).

1. Starten Sie `pgwire-recorder replay` im Hintergrund oder als Container-Dienst.
2. Führen Sie Ihre Tests gegen die Adresse aus `--listen` aus.
3. Beenden Sie das Werkzeug mit `SIGTERM`.
4. Werten Sie den Exit-Code aus (siehe [Exit-Codes](#exit-codes)).

Zeichnen Sie in der Testautomatisierung auf, lassen Sie die Tests vor dem
`SIGTERM` zur Ruhe kommen: Eine Anfrage, die nach Ablauf von
`--shutdown-timeout` noch läuft, fehlt in der Aufzeichnung, und das Werkzeug
endet mit Exit-Code 4, wenn vorher kein anderer Fehler gemerkt wurde (siehe
[Herunterfahren mit Frist](#herunterfahren-mit-frist)).

### Herunterfahren mit Frist

Auf `Strg+C` oder `SIGTERM` fährt das Werkzeug kontrolliert herunter, beim
Aufzeichnen wie beim Wiedergeben. `--shutdown-timeout` (Umgebungsvariable
`PGWIRE_RECORDER_SHUTDOWN_TIMEOUT`) begrenzt, wie lange es dabei auf laufende
Anfragen wartet; die Frist zählt ab dem ersten Signal.

* Der Wert ist `0` oder eine ganze Zahl mit genau einer Einheit `ms`, `s` oder
  `m` in Kleinbuchstaben, etwa `500ms`, `5s` oder `2m`; Standard ist `5s`. `0`
  schaltet die Frist ab, dann wartet das Werkzeug ohne Grenze. Jeder andere Wert
  ist ein ungültiger Aufruf (`PGR-E2001`), auch `5` ohne Einheit, `1.5s`, `5S`
  und `1m30s`. Eine leere Umgebungsvariable gilt als nicht gesetzt; die Option
  geht ihr vor.
* Ist beim Signal keine Verbindung offen, endet das Werkzeug sofort; `record`
  schreibt die Aufzeichnung auch ohne Sitzungen.
* Ein zweites `Strg+C` oder `SIGTERM` lässt die Frist sofort ablaufen, auch beim
  Wert `0`. Jedes weitere Signal bleibt ohne Wirkung: Das Werkzeug beendet die
  Verbindungen, schreibt die Aufzeichnung und endet mit seinem Exit-Code.
* Beim Ablauf beendet das Werkzeug jede noch laufende Verbindung. Eine
  Verbindung ohne begonnene Anfrage endet ohne Meldung; eine mit begonnener,
  nicht abgeschlossener Anfrage erhält die Fehlermeldung `PGR-E4006`, wenn sie
  sie binnen einer Sekunde annimmt, und der Exit-Code ist 4, wenn vorher kein
  anderer Fehler gemerkt wurde; sonst gilt der Exit-Code des ersten gemerkten
  Fehlers, und ein Fehler beim Schreiben der Aufzeichnung ergibt 3 (siehe
  [Exit-Codes](#exit-codes)).
* Die Frist begrenzt nur das Warten auf die Verbindungen. Das Schreiben der
  Aufzeichnung danach bricht sie nicht ab.

Mit `--fail-on-unconsumed` wertet das Werkzeug es als Fehler (`PGR-E5002`,
Exit-Code 5), wenn ein Test nicht alle aufgezeichneten Anfragen ausführt oder
eine aufgezeichnete Sitzung nie verwendet; ohne die Option ist es eine Warnung
(`PGR-W2001`), und der Exit-Code bleibt 0. Die Option nimmt keinen Wert oder
`=true` beziehungsweise `=false`; jeder andere Wert, auch `=` ohne Wert und
`=1`, ist ein ungültiger Aufruf (`PGR-E2001`). Die Umgebungsvariable
`PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` nimmt `true` oder `false`, leer gilt sie als
nicht gesetzt; die Option geht ihr vor. Ein anderer Wert der Variable ist ein
ungültiger Aufruf (`PGR-E2001`), auch wenn Sie die Option zugleich angeben. Die Option gibt es nur bei `replay`.

## 5. Einstellungen

Die Einstellungen lassen sich über Optionen, über Umgebungsvariablen und über eine
Konfigurationsdatei festlegen. Ausnahme sind die benannten Verbindungen (nur in der
Datei). Gilt dieselbe Einstellung mehrfach, setzt sich das Argument vor der
Umgebungsvariablen vor der Konfigurationsdatei vor dem Standardwert durch. Ein leerer
Wert auf der Kommandozeile (`--listen=`) ist ein ungültiger Aufruf (`PGR-E2001`),
auch wenn die Umgebungsvariable gesetzt ist; nur eine leere Umgebungsvariable gilt als
nicht gesetzt.

| Option | Betriebsart | Umgebungsvariable | Standard |
|---|---|---|---|
| `--listen` | `record`, `replay` | `PGWIRE_RECORDER_LISTEN` | Pflicht |
| `--upstream` | `record`, `play` | `PGWIRE_RECORDER_UPSTREAM` | Pflicht (`host:port` oder Name einer Verbindung) |
| `--output` | `record` | `PGWIRE_RECORDER_OUTPUT` | Pflicht |
| `--force` | `record` | `PGWIRE_RECORDER_FORCE` | `false` |
| `--input` | `replay`, `play` | `PGWIRE_RECORDER_INPUT` | Pflicht |
| `--fail-on-unconsumed` | `replay` | `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` | `false` |
| `--shutdown-timeout` | `record`, `replay` | `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` | `5s` (`0` ohne Frist; Einheit `ms`, `s` oder `m`) |
| `--user` | `play` | `PGWIRE_RECORDER_USER` | Benutzer der benutzten Verbindung, ohne ihn der der Aufzeichnung |
| `--database` | `play` | `PGWIRE_RECORDER_DATABASE` | Datenbank der benutzten Verbindung, ohne sie die der Aufzeichnung |
| `--continue-on-error` | `play` | `PGWIRE_RECORDER_CONTINUE_ON_ERROR` | `false` |
| `--allow-recorded-errors` | `play` | `PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS` | `false` |
| `--upstream-tls` | `play` | `PGWIRE_RECORDER_UPSTREAM_TLS` | ohne die Option entscheidet `sslmode` der benutzten Verbindung, ohne beides `false` |
| `--finish-session-on-interrupt` | `play` | `PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT` | `false` |
| `--upstream-ca` | `play` | `PGWIRE_RECORDER_UPSTREAM_CA` | keine Datei, nur der Zertifikatsspeicher des Systems |
| `--config` | `record`, `replay`, `play`, `config show` | `PGWIRE_RECORDER_CONFIG` | `.pgwire-recorder.yaml` im aktuellen Verzeichnis |
| `--log-level` | `record`, `replay`, `play` | `PGWIRE_RECORDER_LOG_LEVEL` | `info` |

Das Passwort, mit dem sich `play` anmeldet, ist keine Option. Schreibt die benutzte
Verbindung einen Passwortteil (siehe [Konfigurationsdatei](#konfigurationsdatei)), gilt
der eingesetzte Wert seines Platzhalters, und die Umgebungsvariable
`PGWIRE_RECORDER_PASSWORD` bleibt unbeachtet. Hat die Verbindung keinen Passwortteil,
oder geben Sie `--upstream` als `host:port` an, gilt der Wert von
`PGWIRE_RECORDER_PASSWORD`; eine leere Variable gilt als nicht gesetzt. Das Passwort
kommt nur zum Einsatz, wenn die Datenbank eines verlangt: Verlangt sie keines
(`trust`), bleibt eine gesetzte Variable ohne Wirkung, und die Anmeldung gelingt. Fehlt
das Passwort, obwohl die Datenbank eines verlangt, endet `play` mit `PGR-E4005`. Keine
Meldung und keine Log-Zeile von `play` und keine Ausgabe von `config show` nennt das
Passwort; `config show` nennt den Namen `PGWIRE_RECORDER_PASSWORD`, wenn die Variable
gesetzt ist, nie ihren Wert. `play` sendet das Passwort unverändert; ein Passwort mit
Zeichen, die SASLprep ändert, kann bei SCRAM scheitern.

### Konfigurationsdatei

Das Werkzeug verwendet genau **eine** Konfigurationsdatei, in dieser Reihenfolge:

1. die Datei aus `--config`,
2. der Pfad aus der Umgebungsvariable `PGWIRE_RECORDER_CONFIG`,
3. die Datei `.pgwire-recorder.yaml` im aktuellen Verzeichnis, falls es sie gibt.

Fehlt eine mit `--config` oder `PGWIRE_RECORDER_CONFIG` genannte Datei, startet das
Werkzeug nicht (`PGR-E2004`); fehlt `.pgwire-recorder.yaml`, liest es keine Datei.

Die Schlüssel heißen wie die Optionen, mit `_` statt `-` und ohne `--`; einen
Schlüssel `config` gibt es nicht, und Groß- und Kleinschreibung zählt. Einstellungen
für alle Betriebsarten (`log_level`) und benannte Verbindungen stehen oben, die
übrigen in einem Abschnitt je Betriebsart (`record:`, `replay:`, `play:`):

```yaml
log_level: info
connections:
  lokal: "postgresql://dev@localhost:5432/myapp"
  ci: "postgresql://${CI_DB_HOST}:5432/myapp"
  test: "postgresql://tester:${DB_PASSWORD}@localhost:5432/myapp"
record:
  upstream: lokal
  listen: 127.0.0.1:15432
  output: ./recordings/users.yaml
play:
  upstream: lokal
  input: ./recordings/users.yaml
```

Eine Verbindung ist eine URL der Form
`postgresql://[benutzer[:passwort]@]host[:port]/datenbank[?sslmode=…]`; die Datenbank
ist Pflicht, ohne Port gilt `5432`, und eine IPv6-Adresse steht in eckigen Klammern
(`postgresql://[::1]:5432/db`). Der Name einer Verbindung enthält weder `:` noch `@`
noch `$`. Danach genügt `--upstream lokal`, ebenso `PGWIRE_RECORDER_UPSTREAM=lokal`
oder der Schlüssel `upstream`. Ein Name gilt nur in genau dieser Schreibweise und nur
aus der gewählten Datei; ohne Datei ist ein Wert ohne `:` ungültig. Jeder Wert von
`--upstream` und `PGWIRE_RECORDER_UPSTREAM` ist der Name einer Verbindung oder hat die
Form `host:port`, auch einer, der wegen der Priorität nicht gilt; sonst startet das
Werkzeug nicht (`PGR-E2001`, im Schlüssel `upstream` `PGR-E2004`).

Bei `play` verbindet das Werkzeug zu Host und Port der URL und meldet jede Sitzung
mit Benutzer und Datenbank der URL an; `--user` und `--database` gehen ihnen vor, und
erst ohne beides gelten Benutzer und Datenbank der Aufzeichnung. Mit dem Beispiel
oben spielt `play` also als `dev` in die Datenbank `myapp` ein, gleich, welche
Datenbank die Aufzeichnung nennt. Mit `--upstream test` meldet sich `play` als
`tester` mit dem Wert von `DB_PASSWORD` an. Die Log-Zeile beim Start nennt auch hier
nur `host:port`.

Bei `record` zählen nur Host und Port der URL; Benutzer, Passwort und Datenbank
vermittelt die Anwendung selbst. Das Werkzeug verbindet zu `host:port`, mit dem Port
so, wie er geschrieben ist, und einem Host mit `:` in eckigen Klammern; die Log-Zeile
beim Start nennt diese Adresse, nie Benutzer, Passwort oder Datenbank.

Der Parameter `sslmode` kennt `disable` (Standard) und `require`. Bei `record` ist eine
Verbindung mit `require` ungültig (`PGR-E2004`), weil `record` unverschlüsselt zur
Datenbank verbindet; eine Verbindung mit `require`, die Sie nicht benutzen, bleibt
gültig. Bei `play` wünscht `require` TLS zum Server, mit Prüfung des Zertifikats; ein
gesetztes `--upstream-tls` geht `sslmode` vor, auch mit `false`. Die Schlüssel
`upstream_tls` und `upstream_ca` stehen im Abschnitt `play:`.

Einen Platzhalter `${VAR}` in der URL einer Verbindung ersetzt das Werkzeug beim Start
aus der gleichnamigen Umgebungsvariable, einmal und nur für die Verbindung, die Sie
benutzen, bei `record` nur in Host und Port. Der Wert steht unverändert an seiner
Stelle, auch mit `@`, `:` oder `/`; eine leere Variable gilt als nicht gesetzt. Ein
Platzhalter darf in jedem Teil der URL stehen außer zwischen eckigen Klammern, im Namen
eines Parameters und in `sslmode`, im Port nur neben Ziffern; außerhalb einer URL ist er
ungültig (`PGR-E2004`). In jedem Wert der Datei steht `$$` für ein `$`, sodass `$${VAR}`
wörtlich `${VAR}` ergibt. Ein Passwort, das nicht genau ein Platzhalter ist, und ein
Parameter `password` sind ein Klartext-Passwort (`PGR-E2006`), in jeder Verbindung der
Datei, auch einer, die Sie nicht benutzen. Eine nicht gesetzte Variable der benutzten
Verbindung ist `PGR-E2005`, ein Port, der nach dem Einsetzen keine Zahl von 1 bis 65535
ist, `PGR-E2004`. Bei `play` ist der eingesetzte Wert eines Platzhalters im Passwort das
Passwort der Anmeldung.

Wahrheitswerte lauten `true` oder `false`, mit oder ohne Anführungszeichen; `True`,
`yes` und `1` sind ungültig. Für jeden Wert gilt dieselbe Form wie für die Option,
eine Dauer also etwa `0` oder `5s`. Ein leerer Wert, `null`, eine Liste oder
Abbildung an der Stelle eines Werts sowie Anker, Aliase, Merge-Schlüssel und
ausdrücklich geschriebene Tags (`!!str`) sind ungültig, ebenso ein Schlüssel, der
in derselben Abbildung zweimal steht, und ein zweites YAML-Dokument in der Datei.
Ein Abschnitt ohne Inhalt (`record:` allein) ist ungültig, `{}` setzt nichts; eine
leere Datei und eine nur mit Kommentaren sind gültig und setzen ebenfalls nichts. Ein relativer Pfad gilt relativ zum
aktuellen Verzeichnis, nicht zum Verzeichnis der Datei.

Die Datei ist eine reguläre Datei in UTF-8; ein BOM am Anfang ist erlaubt, andere
Kodierungen wie UTF-16 sind ungültig. Zeilenenden sind `\n`, `\r\n` oder `\r`. Ein
Verzeichnis, eine Pipe oder ein Gerät als Datei lehnt das Werkzeug ab (`PGR-E2004`),
auch `--config /dev/stdin` mit einer Pipe; eine umgeleitete Datei liest es.

Das Werkzeug prüft beim Start die ganze Datei, auch den Abschnitt einer anderen
Betriebsart, und endet beim ersten Fehler. Die Meldung nennt die Stelle (Schlüssel
oder Verbindung), nie den Wert.

### Log-Ausgaben

Log-Zeilen und Fehlermeldungen gehen nach `stderr`; auf `stdout` stehen nur die
Hilfe, die Ausgabe von `version` und die von `config show`.

`--log-level` (Umgebungsvariable `PGWIRE_RECORDER_LOG_LEVEL`) legt fest, wie viel
das Werkzeug schreibt. Es gibt genau vier Stufen; jede zeigt ihre Zeilen und die
der strengeren:

| Stufe | zeigt |
|---|---|
| `error` | Fehler, die eine Verbindung beenden, bei `play` jeder Fehler nach dem Start |
| `warn` | zusätzlich Warnungen, jede mit ihrem Code |
| `info` (Standard) | zusätzlich Start und Ende von `record`, `replay` und `play`, bei `play` das erste Abbruchsignal, und den Beginn des Herunterfahrens |
| `debug` | zusätzlich Ereignisse je Verbindung; welche, kann sich ändern |

Der Wert wird genau so geschrieben, in Kleinbuchstaben. Jeder andere Wert ist ein
ungültiger Aufruf (`PGR-E2001`), auch ein leerer (`--log-level=`), ein
großgeschriebener (`INFO`) und ein anderer Name (`warning`, `trace`, `off`). Eine
leere Umgebungsvariable gilt als nicht gesetzt; eine mit ungültigem Wert ist
`PGR-E2001`, auch wenn Sie die Option zugleich angeben. Die Option steht nach dem
Kommando; nennen Sie sie mehrfach, gilt die letzte.

Eine Log-Zeile hat das Format `logfmt`: zuerst `time` (Ortszeit nach RFC 3339
mit Millisekunden und Zonenversatz), `level` (`DEBUG`, `INFO`, `WARN`, `ERROR`)
und `msg`, danach weitere Angaben. Ein Fehler trägt `code` und `error` (den
Fehlertext), eine Warnung `code`. Die Zeile beim Beginn des Herunterfahrens
trägt `sessions`, die Zahl der noch offenen Verbindungen, auch `0`. Verlassen
können Sie sich auf `level`, `code`, `error` und `sessions`; der Text von `msg`
und die übrigen Angaben können sich ändern.

```text
time=2026-10-06T14:03:12.481+02:00 level=WARN msg="…" code=PGR-W2001
time=2026-10-06T14:03:12.482+02:00 level=ERROR msg=Fehler code=PGR-E5001 error="Replay [PGR-E5001]: …"
```

Scheitert der Start, schreibt das Werkzeug auf jeder Stufe genau den Fehlertext
als letzte Zeile, ohne `time` und `level`, und vorher nichts anderes; das gilt
auch, wenn der Wert von `--log-level` selbst ungültig ist.

Mit `--help` oder `-h` (ebenso `--h` und `-help`, auch mit `=` und einem Wert)
zeigt das Werkzeug die Hilfe an und endet mit Exit-Code 0, mit
`pgwire-recorder version` die Version. Steht ein Kommando vorn, ist es dessen
Hilfe, sonst die allgemeine. Die Hilfe geht jeder Prüfung vor: Unbekannte
Optionen, ungültige Werte, fehlende Pflichtoptionen und ungültige
Umgebungsvariablen meldet das Werkzeug dann nicht. Auch an der Stelle eines
Optionswerts ist `--help` die Hilfe; als Wert geht es nur mit `=`
(`--input=--help`), und nach `--` ist es ein gewöhnliches Argument. Das erste
`--` beendet die Optionen an jeder Stelle, auch dort, wo ein Optionswert stünde:
`--input --` ist eine Option ohne Wert (`PGR-E2001`), den Wert `--` geben Sie
mit `--input=--` an. `record`, `replay` und `play` nehmen nach `--` kein Argument an
(`PGR-E2001`). Eine Option `--version` gibt es nicht; die Version zeigt
`pgwire-recorder version`.

### Exit-Codes

| Exit-Code | Bedeutung |
|---|---|
| 0 | erfolgreich beendet, auch nach `Strg+C` oder `SIGTERM`, wenn zuvor kein Fehler auftrat |
| 1 | sonstiger Fehler (`PGR-E1000`) |
| 2 | ungültiger Aufruf oder ungültige Konfiguration |
| 3 | Aufzeichnung ungültig oder nicht zugreifbar |
| 4 | Netzwerk- oder Datenbankfehler, auch eine Anfrage, die beim Beenden die Frist `--shutdown-timeout` unvollständig beendet hat (`PGR-E4006`) |
| 5 | Abweichung bei der Wiedergabe, mit `--fail-on-unconsumed` auch nicht gestellte Anfragen oder nie verwendete Sitzungen |
| 6 | nicht unterstützte Funktion des Protokolls |

Ein Fehler, der nur eine Verbindung betrifft, beendet diese Verbindung. Das
Werkzeug läuft weiter und liefert den Exit-Code des ersten aufgetretenen Fehlers
erst, wenn Sie es beenden. Scheitert beim Beenden das Schreiben der
Aufzeichnung, ist der Exit-Code 3, auch wenn vorher ein anderer Fehler auftrat.
Endet eine Verbindung durch einen Fehler und bleiben
dabei Anfragen ihrer Sitzung offen, zählt der Fehler, der sie beendet hat, vor
`PGR-E5002`; beide stehen im Log. Nie verwendete Sitzungen zählen zuletzt und
bestimmen den Exit-Code nur, wenn vorher kein Fehler auftrat. Scheitert schon
der Start, prüft das Werkzeug keine Sitzungen.

## 6. Rollen und Rechte

Das Werkzeug verwaltet keine Benutzer und Rollen und prüft keine Zugangsdaten.
Beim Aufzeichnen leitet es Benutzer und Datenbank der Anwendung an die Datenbank
weiter, ohne Passwort; beim Wiedergeben nimmt es jede Anmeldung an, gleich mit
welchem Benutzer, welcher Datenbank und welchem Passwort; beim Einspielen meldet es
sich mit Benutzer, Datenbank und, wenn die Datenbank eines verlangt, dem Passwort an
der Datenbank an.
Betreiben Sie es nur in einer kontrollierten Testumgebung, und lassen Sie es nur auf der Adresse
lauschen, die Sie mit `--listen` angegeben haben.

## 7. Fehlerbehebung

Jede Meldung des Werkzeugs trägt einen Code der Form `PGR-E…` (Fehler) oder
`PGR-W…` (Warnung). Der Text eines Fehlers ist eine Zeile und beginnt mit der
Fehlerklasse und dem Code, zum Beispiel `Replay [PGR-E5001]: …`; die Klassen
heißen `sonstiger Fehler`, `Konfiguration`, `Recording`, `Netzwerk`, `Replay` und
`nicht unterstützt`. Zeilenumbrüche im Text, auch aus Meldungen des
Betriebssystems, ersetzt das Werkzeug durch ein Leerzeichen. Derselbe Text steht
im Log unter `error` und in der Fehlermeldung, die Ihre Anwendung erhält
(Schweregrad `FATAL`, SQLSTATE `0A000` bei nicht unterstützten Funktionen,
`08006` bei Netzwerkfehlern, sonst `XX000`). Verlassen können Sie sich auf Klasse
und Code; der Text danach kann sich ändern. Entstehen beim Ende einer Verbindung
mehrere Fehler zugleich, etwa der Abbruch der Verbindung und ein Fehler beim
Schreiben der Aufzeichnung, steht jeder in einer eigenen Zeile; der erste zählt
für den Exit-Code.

### Fehlercodes

| Code | Bedeutung | Ursache und Lösung |
|---|---|---|
| `PGR-E1000` | sonstiger Fehler | Ein Defekt des Werkzeugs; keine Eingabe löst ihn aus. Melden Sie das Problem mit der Ausgabe des Laufs. |
| `PGR-E2001` | ungültiger Aufruf | Eine Option fehlt, ist unbekannt oder hat einen ungültigen Wert, oder `--upstream-ca` steht, obwohl `play` kein TLS nutzt (ohne `--upstream-tls` und ohne `sslmode=require` der Verbindung, oder mit `--upstream-tls=false`). Prüfen Sie den Aufruf mit `--help`. |
| `PGR-E2002` | Zieldatei existiert bereits | Wählen Sie einen anderen Dateinamen, oder ergänzen Sie `--force`, um die Datei zu ersetzen. |
| `PGR-E2004` | Konfigurationsdatei nicht lesbar oder ungültig | Die Meldung nennt den Schlüssel oder die Verbindung. Prüfen Sie YAML, Schlüssel, Abschnitt, Werte und `sslmode` (erlaubt sind `disable` und `require`). |
| `PGR-E2005` | Umgebungsvariable eines Platzhalters nicht gesetzt | Die Meldung nennt die Verbindung und die erste fehlende Variable. Setzen Sie sie mit einem nicht leeren Wert; eine leere Variable gilt als nicht gesetzt. Bei `record` zählen nur die Variablen in Host und Port. |
| `PGR-E2006` | Klartext-Passwort in der Konfigurationsdatei | Die Meldung nennt die Verbindung. Ersetzen Sie das Passwort in der URL durch genau einen Platzhalter `${VAR}`, und entfernen Sie einen Parameter `password`. Das gilt für jede Verbindung der Datei, auch eine, die Sie nicht benutzen. |
| `PGR-E2007` | Datei von `--upstream-ca` nicht verwendbar | Die Datei fehlt, ist nicht lesbar, ist keine reguläre Datei oder enthält keine Zertifikate im PEM-Format. Die Meldung nennt die Option und den Grund, nicht den Pfad; der Exit-Code ist 2. Prüfen Sie Pfad, Rechte und Inhalt der Datei. |
| `PGR-E3001` | Aufzeichnung nicht lesbar oder nicht schreibbar | Die Datei fehlt, Sie haben keine Rechte, `--output` ist ein Verzeichnis, oder das Verzeichnis von `--output` fehlt. Prüfen Sie Pfad und Dateirechte. |
| `PGR-E3002` | unbekannte Version der Aufzeichnung | Die Datei stammt aus einer anderen Programmversion. Zeichnen Sie mit der verwendeten Version erneut auf. |
| `PGR-E3003` | Aufzeichnung beschädigt | Die Datei ist unvollständig oder verändert. Zeichnen Sie erneut auf. |
| `PGR-E3004` | Aufzeichnung ohne verwendbare Sitzung | `replay` meldet den Fehler beim Start. Beim Aufzeichnen hat keine Verbindung eine Anfrage gestellt, oder die Verbindungen haben nur Lebendprüfungen gesendet (Anfragen nur aus Leerraum und Kommentaren, siehe [Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)). Zeichnen Sie einen Ablauf mit mindestens einer anderen Anfrage auf. |
| `PGR-E4000` | Verbindung nicht anzunehmen | Das Werkzeug konnte eine eingehende Verbindung nicht annehmen, zum Beispiel weil zu viele Dateien offen sind. Es nimmt danach weiter Verbindungen an. Prüfen Sie die Grenzen des Systems. |
| `PGR-E4003` | Verbindung unerwartet beendet | Die Verbindung brach mitten in einer Anfrage ab. Prüfen Sie Netzwerk, Datenbank und Anwendung. |
| `PGR-E4001` | Adresse nicht nutzbar | Der Port aus `--listen` ist belegt oder nicht erlaubt. Wählen Sie einen freien Port. |
| `PGR-E4002` | Datenbank nicht erreichbar | Prüfen Sie `--upstream`, die Datenbank und das Netzwerk. Bei `play` mit TLS gehört dazu: Die Antwort auf die Anfrage nach Verschlüsselung ist weder `S` noch `N`, oder die Verbindung endet davor. |
| `PGR-E4004` | Datenbank beantwortet eine eingespielte Anfrage mit einem Fehler | Die Meldung nennt die Sitzung und die Nummer der Anfrage in der Aufzeichnung, dazu SQLSTATE und Meldung der Datenbank, nicht den Text der Anfrage. Prüfen Sie Benutzer, Rechte und den Zustand der Datenbank, oder starten Sie mit `--continue-on-error`. |
| `PGR-E4005` | Anmeldung oder TLS zur Datenbank abgelehnt oder nicht unterstützt | Beim Einspielen hat die Datenbank die Anmeldung abgelehnt, etwa weil der Benutzer fehlt oder das Passwort falsch ist, oder sie verlangt ein Passwort, das `play` nicht hat, oder ein Anmeldeverfahren, das `play` nicht kann. Die Meldung nennt die Sitzung; bei einer Ablehnung außerdem SQLSTATE und Meldung der Datenbank, bei fehlendem Passwort, dass `play` keines hat. `play` meldet sich mit Klartext-Passwort, MD5 und SCRAM-SHA-256 an; bietet die Datenbank nur `SCRAM-SHA-256-PLUS` an oder verlangt sie Kerberos, GSSAPI oder SSPI, endet `play` ebenfalls mit diesem Code. Dasselbe gilt für eine Datenbank, die Verbindungen ohne Verschlüsselung ablehnt, wenn `play` kein TLS wünscht. Mit TLS endet `play` mit diesem Code auch, wenn die Datenbank TLS ablehnt (Antwort `N`, ohne Rückfall auf eine Verbindung ohne TLS) oder ihr Zertifikat nicht angenommen wird: Das Zertifikat lässt sich weder mit dem Speicher des Systems noch mit `--upstream-ca` überprüfen, der Host steht nicht im Zertifikat, oder das Zertifikat ist abgelaufen. Die Meldung nennt die Sitzung und den Grund, nie Pfad oder Inhalt der Datei aus `--upstream-ca`. Prüfen Sie Benutzer und Datenbank (`--user`, `--database`), das Passwort (siehe [Einstellungen](#5-einstellungen)), den Host der Verbindung und `--upstream-ca`. |
| `PGR-E4006` | Anfrage beim Beenden unvollständig | Die Frist `--shutdown-timeout` ist abgelaufen, oder ein zweites Signal hat sie ablaufen lassen, während eine Anfrage lief. Die Meldung nennt die Sitzung und die Anfrage; beim Aufzeichnen fehlt diese Anfrage in der Aufzeichnung. Lassen Sie die Anwendung vor dem Beenden zur Ruhe kommen, oder wählen Sie eine längere Frist (siehe [Herunterfahren mit Frist](#herunterfahren-mit-frist)). |
| `PGR-E5001` | Abweichung bei der Wiedergabe | Ihre Anwendung hat eine andere Anfrage gestellt als aufgezeichnet. Die Meldung nennt die erwartete und die empfangene Anfrage. Zeichnen Sie erneut auf, oder korrigieren Sie die Anwendung. |
| `PGR-E5002` | aufgezeichnete Anfragen oder Sitzungen nicht verbraucht | Ihr Test hat weniger Anfragen gestellt oder weniger Verbindungen geöffnet als aufgezeichnet, und `--fail-on-unconsumed` ist gesetzt. |
| `PGR-E5003` | Anfrage ohne aufgezeichnete Sitzung | Ihre Anwendung hat auf mehr Verbindungen Anfragen gestellt, als Sitzungen aufgezeichnet sind. Zeichnen Sie den Ablauf erneut auf, oder öffnen Sie weniger Verbindungen. |
| `PGR-E6001` | nicht unterstützte Nachricht oder Funktion | Die Anwendung nutzt eine Funktion, die das Werkzeug nicht unterstützt, zum Beispiel `COPY`; beim Aufzeichnen verlangt die Datenbank ein Passwort; oder die Datenbank antwortet bei `play` mit einem COPY-Datenstrom. Verwenden Sie diese Funktion im aufgezeichneten Ablauf nicht; verlangt die Datenbank beim Aufzeichnen ein Passwort, lassen Sie die Anmeldung dort ohne Passwort zu. |
| `PGR-E6002` | nicht unterstützte Protokollversion | Das Werkzeug unterstützt Version 3.0 des PostgreSQL-Protokolls. Verwenden Sie einen Treiber, der sie nutzt. |

### Warnungen

| Code | Bedeutung | Hinweis |
|---|---|---|
| `PGR-W2001` | Wiedergabe endete vor der letzten aufgezeichneten Anfrage, oder aufgezeichnete Sitzungen wurden nie verwendet | Ihr Test hat nicht alle aufgezeichneten Anfragen ausgeführt oder weniger Verbindungen mit Anfragen geöffnet als aufgezeichnet. Mit `--fail-on-unconsumed` wird das zum Fehler (`PGR-E5002`). |
| `PGR-W3001` | Abbruchwunsch nicht weitergeleitet | Die Anwendung hat versucht, eine laufende Anfrage abzubrechen. Das Werkzeug leitet diesen Wunsch nicht weiter und schließt die Verbindung. |
| `PGR-W3003` | Verbindung ohne PostgreSQL-Protokoll | Etwas anderes als ein PostgreSQL-Client hat den Port angesprochen, zum Beispiel ein HTTP-Gesundheitscheck. Das Werkzeug schließt die Verbindung; der Exit-Code ändert sich nicht. Ein reiner TCP-Check ohne Daten erzeugt keine Warnung. |

### Die Anwendung kann sich nicht verbinden

Fordert Ihr Treiber eine verschlüsselte Verbindung an und bricht ab, wenn sie
abgelehnt wird (etwa mit `sslmode=require`), schalten Sie die Verschlüsselung in der
Verbindungszeichenfolge ab (`sslmode=disable`), oder lassen Sie sie mit
`sslmode=prefer` unverschlüsselt verbinden. Erhält Ihre Anwendung beim Aufzeichnen
`PGR-E6001` mit dem Hinweis auf ein Anmeldeverfahren, verlangt die Datenbank ein
Passwort; lassen Sie die Anmeldung dort ohne Passwort zu.

## 8. FAQ

**Brauche ich für die Wiedergabe eine Datenbank?**
Nein. Die Wiedergabe beantwortet alle aufgezeichneten Anfragen aus der Datei.

**Kann ich eine Aufzeichnung in die Versionsverwaltung legen?**
Ja. Die Aufzeichnung ist Text (YAML) und ändert sich nur, wenn sich die
Aufzeichnung ändert. Beachten Sie den Hinweis zu vertraulichen Daten in der Einleitung.

**Kann ich die Datei von Hand bearbeiten?**
Das Werkzeug unterstützt das nicht. Zeichnen Sie stattdessen erneut auf.

**Funktioniert die Wiedergabe, wenn ich die Anfragen umformuliere?**
Nein. Die Anfrage muss Zeichen für Zeichen der aufgezeichneten entsprechen,
auch in Leerzeichen, Kommentaren und Groß- und Kleinschreibung. Beim erweiterten
Protokoll gilt das auch für Namen, Parametertypen, Parameterwerte, Formate und
das Zeilenlimit. Ausgenommen sind nur Anfragen, die nur aus Leerraum und
Kommentaren bestehen (Lebendprüfungen, siehe
[Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)).

**Kann ich Aufzeichnungen auf einem anderen Rechner verwenden?**
Ja. Die Datei enthält keine rechnerspezifischen Angaben des Werkzeugs.

## 9. Glossar

| Begriff | Bedeutung |
|---|---|
| Aufzeichnung | Datei mit der aufgezeichneten Kommunikation zwischen Anwendung und Datenbank |
| Verbindung | eine Datenbankverbindung Ihrer Anwendung; in der Aufzeichnung heißt sie Sitzung |
| Einfaches Anfrageprotokoll | PostgreSQL-Protokollvariante, in der jede Anfrage als einzelner Text gesendet wird |
| Erweitertes Anfrageprotokoll | Protokollvariante mit vorbereiteten Anweisungen und Parametern |
| Upstream | die echte Datenbank, an die `record` weiterleitet |
| Meldungscode | Kennung einer Meldung des Werkzeugs, zum Beispiel `PGR-E5001` |

## 10. Support und Kontakt

Melden Sie Fehler und Fragen über das Projekt-Repository. Geben Sie die
Programmversion, den Aufruf (ohne Zugangsdaten) und die Ausgabe mit
`--log-level debug` an. Entfernen Sie vertrauliche Daten, bevor Sie Ausgaben
oder Aufzeichnungen weitergeben.

## 11. Änderungshistorie

Es gibt keine veröffentlichte Version.
