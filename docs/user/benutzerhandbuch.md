# Benutzerhandbuch: pgwire-recorder

Version: 0.1  
Software-Version: noch nicht veröffentlicht  
Stand: 04.10.2026  
Autor: Projektteam pgwire-recorder  
Gültigkeitsbereich: gilt für `pgwire-recorder` ab der ersten veröffentlichten Version

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
* Zum Aufzeichnen: eine erreichbare PostgreSQL-Datenbank.
* Linux, macOS oder Windows auf `amd64` oder `arm64`, oder ein
  Container-Laufzeitsystem.
* Die Verbindung zum Werkzeug läuft ohne Verschlüsselung, es sei denn, Sie
  stellen Zertifikat und Schlüssel bereit (siehe
  [Verschlüsselte Verbindungen annehmen](#verschlüsselte-verbindungen-annehmen) und
  [Fehlerbehebung](#7-fehlerbehebung)).

### Wichtig: Aufzeichnungen können vertrauliche Daten enthalten

Eine Aufzeichnung enthält die SQL-Anfragen und die zurückgegebenen Datensätze.
Das Werkzeug erkennt oder maskiert vertrauliche Inhalte nicht. Verwenden Sie
zum Aufzeichnen geeignete Testdaten, und speichern, übertragen und versionieren
Sie Aufzeichnungen so, wie es zu ihrem Inhalt passt.

## 2. Installation

### Binary

Das Binary `pgwire-recorder` gibt es für Linux, macOS und Windows, jeweils für
`amd64` und `arm64`.

1. Laden Sie das Binary für Ihre Plattform von den Releases des Projekt-Repositorys
   auf GitHub herunter, und legen Sie es in ein Verzeichnis Ihres Suchpfads.
2. Machen Sie die Datei unter Linux und macOS ausführbar.
3. Prüfen Sie die Installation mit `pgwire-recorder version`.

### Homebrew (macOS und Linux)

Das Werkzeug liegt in einem eigenen Tap, nicht im Standard-Repository von
Homebrew; `brew install pgwire-recorder` allein findet es deshalb nicht. Sie
brauchen zwei Schritte:

```bash
brew tap pt9912/pgwire-recorder
brew install pgwire-recorder
```

Verlangt Ihre Homebrew-Version, einen Drittanbieter-Tap vor der Installation als
vertrauenswürdig zu markieren, führen Sie diesen Schritt vor der Installation aus;
die Meldung von Homebrew nennt den Befehl.

Der Tap enthält nur veröffentlichte, stabile Versionen. Prüfen Sie die
Installation mit `pgwire-recorder version`.

### Container

Das Docker/OCI-Image für `linux/amd64` und `linux/arm64` enthält das Binary und
braucht keine weitere Software; es läuft mit jedem OCI-kompatiblen
Container-Laufzeitsystem, zum Beispiel Docker oder Podman. Sie finden es in der
GitHub Container Registry (`ghcr.io/pt9912/pgwire-recorder`) und auf Docker Hub
(`pt9912/pgwire-recorder`). Ein Beispiel für eine Wiedergabe in Docker Compose
(setzen Sie die Version ein):

```yaml
services:
  recorder:
    image: ghcr.io/pt9912/pgwire-recorder:<Version>
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

* Mit `--format sqlite` speichert das Werkzeug die Aufzeichnung als SQLite-Datei
  statt als Textdatei. Die SQLite-Datei wächst mit jeder beendeten Sitzung und ist
  für große Aufzeichnungen geeignet; sie ist binär und lässt sich nicht sinnvoll
  vergleichen. Beim Wiedergeben und Einspielen erkennt das Werkzeug das Format
  selbst.
* Mit `--record-timing` hält die Aufzeichnung den zeitlichen Abstand der Anfragen
  fest. Zwei Aufzeichnungen derselben Anwendung unterscheiden sich dann in diesen
  Angaben.
* Existiert die Zieldatei bereits, bricht das Werkzeug ab (`PGR-E2002`). Wollen
  Sie sie ersetzen, ergänzen Sie `--force`.
* Jede Verbindung Ihrer Anwendung mit mindestens einer Anfrage wird als eigene
  Sitzung aufgezeichnet. Verbindungen ohne Anfrage, zum Beispiel
  Probe-Verbindungen eines Connection-Pools, werden nur mit
  `--record-empty-sessions` aufgezeichnet.

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
  Antwort (`PGR-E5001`).
* Die erste Verbindung mit einer Anfrage erhält die erste aufgezeichnete Sitzung,
  die zweite die zweite, und so weiter. Verbindungen ohne Anfrage zählen nicht.
  Mit `--session-assignment connection` zählt stattdessen die Reihenfolge der
  Verbindungen, auch ohne Anfrage; das geht nur mit Aufzeichnungen, die mit
  `--record-empty-sessions` erzeugt wurden (sonst `PGR-E2001`).
  Nutzen Sie die Verbindungen nacheinander; bei gleichzeitiger Nutzung ist die
  Zuordnung nicht festgelegt. Eine Anfrage über die aufgezeichneten Sitzungen
  hinaus wird als Abweichung gemeldet (`PGR-E5003`).

### Verschlüsselte Verbindungen annehmen

Damit nehmen `record` und `replay` verschlüsselte Verbindungen von Anwendungen an,
zum Beispiel wenn ein Treiber Verschlüsselung verlangt.

#### Voraussetzung

Ein Zertifikat samt privatem Schlüssel liegt als PEM-Datei vor. Der Schlüssel darf
nicht mit einem Passwort geschützt sein. Für einen Test erzeugen Sie beides zum
Beispiel mit OpenSSL (der Name `localhost` ist der, den Ihre Anwendung verwendet):

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 30 \
  -keyout server-key.pem -out server.pem \
  -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost"
```

#### Vorgehen

1. Starten Sie `pgwire-recorder record` oder `pgwire-recorder replay` mit den
   Optionen `--tls-cert` und `--tls-key`:

   ```bash
   pgwire-recorder replay \
     --listen 0.0.0.0:5432 \
     --input ./recordings/users.yaml \
     --tls-cert ./certs/server.pem \
     --tls-key ./certs/server-key.pem
   ```

2. Verbinden Sie Ihre Anwendung wie gewohnt, mit eingeschalteter Verschlüsselung.

#### Ergebnis

Ihre Anwendung verbindet sich verschlüsselt; Aufzeichnung und Wiedergabe laufen wie
ohne Verschlüsselung. Die Aufzeichnung enthält keine Angabe zur Verschlüsselung und
lässt sich mit und ohne sie wiedergeben.

#### Hinweise

* Die Optionen gelten nur zusammen. Ist eine Datei nicht lesbar oder kein gültiges
  PEM, ist das Zertifikat abgelaufen, ist der Schlüssel mit einem Passwort geschützt
  oder gehören Zertifikat und Schlüssel nicht zusammen, endet der Start mit
  `PGR-E2007`. Relative Pfade gelten ab dem aktuellen Verzeichnis.
* Ob der Name im Zertifikat zum Host passt und ob das Zertifikat vertrauenswürdig
  ist, prüft Ihre Anwendung. Bei einem selbst erzeugten Zertifikat hinterlegen Sie
  es dort als vertrauenswürdig.
* Ist Verschlüsselung eingerichtet, weist das Werkzeug unverschlüsselte
  Verbindungen ab (`PGR-E6003`). Mit `--allow-plaintext` lässt es beides zu.
* Scheitert die Verschlüsselung einer einzelnen Verbindung, schließt das Werkzeug
  sie und warnt (`PGR-W3002`); der Exit-Code ändert sich dadurch nicht.
* Zertifikate der Anwendung (Client-Zertifikate) prüft das Werkzeug nicht. Die
  Verbindung des Werkzeugs zur Datenbank beim Aufzeichnen bleibt unverschlüsselt.
* Der Schlüssel erscheint weder in Meldungen noch in der Anzeige der
  Konfiguration; schützen Sie die Datei selbst mit Dateirechten.

### Eine Aufzeichnung in eine Datenbank einspielen

Damit führen Sie die aufgezeichneten Anfragen erneut aus, zum Beispiel um eine
Komponente zu testen, die Änderungen der Datenbank verarbeitet (Change Data Capture).

#### Voraussetzung

Eine Aufzeichnung liegt vor, und die Zieldatenbank ist erreichbar. Der Benutzer
und die Datenbank aus der Aufzeichnung existieren dort oder Sie geben sie
ausdrücklich an.

#### Vorgehen

1. Setzen Sie bei Bedarf das Passwort in der Umgebungsvariable
   `PGWIRE_RECORDER_PASSWORD` oder legen Sie es als Platzhalter `${VAR}` in einer
   benannten Verbindung der Konfigurationsdatei ab; es gibt dafür keine Option.
2. Starten Sie das Einspielen:

   ```bash
   pgwire-recorder play \
     --upstream postgres:5432 \
     --input ./recordings/users.yaml
   ```

3. Warten Sie, bis das Werkzeug endet.

#### Ergebnis

Die Anfragen der Aufzeichnung sind in der aufgezeichneten Reihenfolge gegen die
Datenbank ausgeführt; sie enthält danach deren Wirkung. Jede Sitzung der
Aufzeichnung läuft über eine eigene Verbindung, die Sitzungen nacheinander.

#### Hinweise

* Das Einspielen läuft nacheinander. Ohne weitere Option gibt es keine
  Wartezeiten. Mit `--keep-timing` stellt das Werkzeug den aufgezeichneten
  zeitlichen Abstand her, sofern Sie die Aufzeichnung mit `--record-timing`
  erzeugt haben. Mit `--timing-mode relative` (Standard) ist der Abstand zur
  vorigen Anfrage nie kürzer als aufgezeichnet. Mit `--timing-mode absolute` liegt
  keine Anfrage früher als im aufgezeichneten Abstand zu einem Bezugspunkt, den
  `--timing-reference` wählt (`connect`: Beginn des Verbindungsaufbaus,
  `first-request`: erste Anfrage); eine Verspätung wird dabei aufgeholt.
  `--timing-mode` und `--timing-reference` brauchen `--keep-timing`.
* Ohne weitere Option vergleicht das Werkzeug die Antworten der Datenbank nicht mit
  der Aufzeichnung. Mit `--compare-responses` prüft es nach jeder Anfrage die
  Struktur der Antwort: die Art und Reihenfolge der Nachrichten, die Spalten (Anzahl,
  Name, Typ), den Befehl, Fehler (Fehlercode) und den Transaktionsstatus. Zeilenwerte,
  Zeilenzahlen und Hinweise der Datenbank vergleicht es nicht. Bei einer Abweichung
  endet das Einspielen mit `PGR-E5004` und Exit-Code 5; mit `--continue-on-error`
  läuft es weiter und endet am Ende mit Exit-Code 5. Ein Fehler der Datenbank, den
  die Aufzeichnung genauso enthält, gilt dann als erwartet und bricht nicht ab; ein
  Fehler, den sie nicht enthält, ist eine Abweichung (`PGR-E5004` statt
  `PGR-E4004`). `--allow-recorded-errors` hat mit Vergleich keine Wirkung; die
  Kombination ist kein Fehler. Reicht die Aufzeichnung einer Anfrage nicht bis zum
  Ende der Antwort, vergleicht das Werkzeug, soweit sie reicht; was der Server danach
  sendet, ist keine Abweichung, außer einer Fehlerantwort, die die Aufzeichnung nicht
  enthält. Endet die Aufzeichnung mit einem Verbindungsende und der Server
  verursacht es genauso, gilt es als erwartet, und das Werkzeug macht mit der
  nächsten Sitzung weiter. Beendet der Server die Verbindung mit einem anderen oder
  ohne aufgezeichneten Fehler, oder antwortet er dort normal, ist das eine
  Abweichung; ein Verbindungsverlust ohne Fehlerantwort (zum Beispiel ein
  Netzwerkabbruch) bleibt `PGR-E4003`. Das Werkzeug meldet je Anfrage die erste
  Abweichung. Bei jedem Ende des Laufs steht eine Zeile im Log (Stufe `info`, bei
  `--log-level warn` also nicht sichtbar) mit der Zahl der eingespielten,
  verglichenen und abweichenden Anfragen. Ein Verbindungsverlust beendet das
  Einspielen immer mit Exit-Code 4, auch nach einer Abweichung.
* Antwortet die Datenbank auf eine Anfrage mit einem Fehler, bricht das Einspielen
  ab (`PGR-E4004`; mit `--compare-responses` gilt stattdessen der Vergleich). Mit `--continue-on-error` läuft es weiter und endet am Ende
  mit Exit-Code 4. Mit `--allow-recorded-errors` gilt ein Fehler nicht, wenn auch
  die aufgezeichnete Anfrage mit einem Fehler beantwortet wurde.
* Mit `--upstream-tls` verbindet sich das Werkzeug verschlüsselt mit der
  Datenbank und prüft deren Zertifikat; ein Überspringen der Prüfung gibt es nicht.
  Verlangt die Datenbank Verschlüsselung und die Option fehlt, oder schlägt die
  Anmeldung fehl, meldet es `PGR-E4005`.
* Trägt die Datenbank ein Zertifikat einer eigenen Zertifizierungsstelle, geben Sie
  deren Zertifikat mit `--upstream-ca` (PEM-Datei) an. Es ergänzt die Zertifikate
  Ihres Systems; die Prüfung bleibt vollständig. Die Option gilt nur mit
  Verschlüsselung; eine nicht lesbare Datei meldet `PGR-E2007`.
* Bei `Strg+C` oder `SIGTERM` endet das Einspielen nach der laufenden Anfrage; mit
  `--finish-session-on-interrupt` erst nach der laufenden Sitzung.
* Das Einspielen führt Anfragen aus. Verwenden Sie es nicht gegen eine Datenbank,
  deren Inhalt Sie nicht verändern dürfen.

### Die gewählte Konfigurationsdatei anzeigen

#### Voraussetzung

Eine Konfigurationsdatei liegt vor, zum Beispiel `.pgwire-recorder.yaml` im aktuellen
Verzeichnis.

#### Vorgehen

1. Führen Sie `pgwire-recorder config show` aus, bei Bedarf mit `--config <Datei>`.

#### Ergebnis

Das Werkzeug nennt die gewählte Datei und gibt ihren Inhalt als eingerückten Baum
aus; Platzhalter `${VAR}` erscheinen unaufgelöst, und aktive
`PGWIRE_RECORDER_*`-Umgebungsvariablen stehen am Ende mit ihrem Namen, ohne Wert.
Der Befehl endet mit Exit-Code 0. Findet er keine Datei, sagt er das und endet
ebenfalls mit Exit-Code 0.

#### Hinweise

* Die Ausgabe enthält nie einen aufgelösten Wert, aber Hosts, Benutzer und Pfade der
  Datei. Prüfen Sie sie, bevor Sie sie weitergeben.
* Ist die Datei ungültig, zeigt das Werkzeug nichts und meldet `PGR-E2004` bis
  `PGR-E2006`.

### Mit einem Datenbanktreiber arbeiten

Das Werkzeug unterstützt das einfache und das erweiterte Anfrageprotokoll von
PostgreSQL. Viele Treiber nutzen standardmäßig das erweiterte Protokoll mit
vorbereiteten Anweisungen. Dafür müssen Sie am Treiber nur Host und Port ändern.

#### Voraussetzung

Ihr Treiber verwendet Version 3.0 des PostgreSQL-Protokolls und lässt sich ohne
Verschlüsselung betreiben, oder Sie stellen dem Werkzeug Zertifikat und Schlüssel
bereit (siehe
[Verschlüsselte Verbindungen annehmen](#verschlüsselte-verbindungen-annehmen)).

#### Vorgehen

1. Starten Sie `pgwire-recorder record` oder `pgwire-recorder replay`.
2. Tragen Sie in der Verbindungszeichenfolge Ihres Treibers Host und Port aus
   `--listen` ein.
3. Stellen Sie die Verschlüsselung in der Verbindungszeichenfolge ab, zum Beispiel
   mit `sslmode=disable`, oder starten Sie das Werkzeug mit Zertifikat und
   Schlüssel.
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
* Sendet ein Treiber nach einer Anforderung `Flush` weitere Nachrichten, ohne auf
  die Antwort zu warten, kann sich die Aufzeichnung zwischen zwei Läufen
  unterscheiden; die Wiedergabe einer vorhandenen Aufzeichnung bleibt davon
  unberührt.

### Das Werkzeug in der Testautomatisierung verwenden

Das Werkzeug fragt nie nach Eingaben und lässt sich über Optionen und
Umgebungsvariablen vollständig steuern (siehe [Einstellungen](#5-einstellungen)).

1. Starten Sie `pgwire-recorder replay` im Hintergrund oder als Container-Dienst.
2. Führen Sie Ihre Tests gegen die Adresse aus `--listen` aus.
3. Beenden Sie das Werkzeug mit `SIGTERM`.
4. Werten Sie den Exit-Code aus (siehe [Exit-Codes](#exit-codes)).

Mit `--fail-on-unconsumed` wertet das Werkzeug es als Fehler, wenn ein Test
nicht alle aufgezeichneten Anfragen ausführt.

## 5. Einstellungen

Die Einstellungen lassen sich über Optionen, über Umgebungsvariablen und über eine
Konfigurationsdatei festlegen. Ausnahmen sind das Passwort (nur über die Umgebung
oder einen Platzhalter in der Datei) und die benannten Verbindungen (nur in der
Datei). Gilt dieselbe Einstellung mehrfach, setzt sich das Argument vor der
Umgebungsvariablen vor der Konfigurationsdatei vor dem Standardwert durch.

| Option | Betriebsart | Umgebungsvariable | Standard |
|---|---|---|---|
| `--listen` | `record`, `replay` | `PGWIRE_RECORDER_LISTEN` | Pflicht |
| `--upstream` | `record`, `play` | `PGWIRE_RECORDER_UPSTREAM` | Pflicht (`host:port` oder Name einer Verbindung) |
| `--output` | `record` | `PGWIRE_RECORDER_OUTPUT` | Pflicht |
| `--force` | `record` | `PGWIRE_RECORDER_FORCE` | `false` |
| `--format` | `record` | `PGWIRE_RECORDER_FORMAT` | `yaml` (oder `sqlite`) |
| `--record-timing` | `record` | `PGWIRE_RECORDER_RECORD_TIMING` | `false` |
| `--record-empty-sessions` | `record` | `PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS` | `false` |
| `--input` | `replay`, `play` | `PGWIRE_RECORDER_INPUT` | Pflicht |
| `--fail-on-unconsumed` | `replay` | `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` | `false` |
| `--session-assignment` | `replay` | `PGWIRE_RECORDER_SESSION_ASSIGNMENT` | `first-request` (oder `connection`) |
| `--user` | `play` | `PGWIRE_RECORDER_USER` | Daten aus der Aufzeichnung |
| `--database` | `play` | `PGWIRE_RECORDER_DATABASE` | Daten aus der Aufzeichnung |
| `--continue-on-error` | `play` | `PGWIRE_RECORDER_CONTINUE_ON_ERROR` | `false` |
| `--allow-recorded-errors` | `play` | `PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS` | `false` |
| `--upstream-tls` | `play` | `PGWIRE_RECORDER_UPSTREAM_TLS` | `false` |
| `--upstream-ca` | `play` | `PGWIRE_RECORDER_UPSTREAM_CA` | — (PEM-Datei, nur mit Verschlüsselung) |
| `--compare-responses` | `play` | `PGWIRE_RECORDER_COMPARE_RESPONSES` | `false` |
| `--finish-session-on-interrupt` | `play` | `PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT` | `false` |
| `--keep-timing` | `play` | `PGWIRE_RECORDER_KEEP_TIMING` | `false` |
| `--timing-mode` | `play` | `PGWIRE_RECORDER_TIMING_MODE` | `relative` (oder `absolute`) |
| `--timing-reference` | `play` | `PGWIRE_RECORDER_TIMING_REFERENCE` | `connect` (oder `first-request`) |
| `--tls-cert` | `record`, `replay` | `PGWIRE_RECORDER_TLS_CERT` | — (PEM-Datei, verlangt `--tls-key`) |
| `--tls-key` | `record`, `replay` | `PGWIRE_RECORDER_TLS_KEY` | — (PEM-Datei, verlangt `--tls-cert`) |
| `--allow-plaintext` | `record`, `replay` | `PGWIRE_RECORDER_ALLOW_PLAINTEXT` | `false` (verlangt `--tls-cert`) |
| — (nur Umgebung) | `play` | `PGWIRE_RECORDER_PASSWORD` | — |
| `--config` | `record`, `replay`, `play`, `config show` | `PGWIRE_RECORDER_CONFIG` | `.pgwire-recorder.yaml` im aktuellen Verzeichnis |
| `--log-level` | `record`, `replay`, `play` | `PGWIRE_RECORDER_LOG_LEVEL` | `info` |

### Konfigurationsdatei

Das Werkzeug verwendet genau **eine** Konfigurationsdatei, in dieser Reihenfolge:

1. die Datei aus `--config`,
2. der Pfad aus der Umgebungsvariable `PGWIRE_RECORDER_CONFIG`,
3. die Datei `.pgwire-recorder.yaml` im aktuellen Verzeichnis, falls es sie gibt.

Die Schlüssel heißen wie die Optionen, mit `_` statt `-` und ohne `--`. Einstellungen
für alle Betriebsarten (`log_level`) und benannte Verbindungen stehen oben, die
übrigen in einem Abschnitt je Betriebsart (`record:`, `replay:`, `play:`):

```yaml
log_level: info
connections:
  lokal: "postgresql://dev@localhost:5432/myapp"
  staging: "postgresql://app:${STAGING_PASSWORD}@staging.example.com:5432/myapp?sslmode=require"
play:
  upstream: staging
  input: ./recordings/users.yaml
  keep_timing: true
  timing_mode: relative
```

Danach genügt `--upstream staging`; ein Name hat Vorrang vor `host:port`. Der
Parameter `sslmode` kennt `disable` (Standard) und `require`. `require` verbindet
verschlüsselt und prüft das Zertifikat der Datenbank; ein ausdrücklich gesetztes
`--upstream-tls` geht dem `sslmode` vor. Bei `record` zählen nur Host und Port der
URL; Benutzer, Passwort und Datenbank vermittelt die Anwendung selbst.

Passwörter geben Sie als Platzhalter `${VAR}` in der URL einer Verbindung an; das
Werkzeug ersetzt ihn aus der gleichnamigen Umgebungsvariable (`$${VAR}` bleibt
wörtlich) und nur für die Verbindung, die Sie benutzen. Fehlt in einer Verbindung
das Passwort, gilt `PGWIRE_RECORDER_PASSWORD`. Ein Klartext-Passwort in der Datei
lehnt das Werkzeug ab (`PGR-E2006`), ebenso eine nicht gesetzte Variable
(`PGR-E2005`) und eine ungültige Datei (`PGR-E2004`).

Wahrheitswerte lauten `true` oder `false`. Mögliche Log-Level sind `error`,
`warn`, `info` und `debug`. Meldungen gehen nach `stderr`.

Mit `--help` oder `-h` zeigt das Werkzeug die Hilfe an, mit
`pgwire-recorder version` die Version.

### Exit-Codes

| Exit-Code | Bedeutung |
|---|---|
| 0 | erfolgreich beendet, auch nach `Strg+C` oder `SIGTERM`, wenn zuvor kein Fehler auftrat |
| 1 | sonstiger Fehler |
| 2 | ungültiger Aufruf oder ungültige Konfiguration |
| 3 | Aufzeichnung ungültig oder nicht zugreifbar |
| 4 | Netzwerk- oder Datenbankfehler |
| 5 | Abweichung bei der Wiedergabe, oder beim Einspielen mit `--compare-responses` eine abweichende Antwort |
| 6 | nicht unterstützte Funktion des Protokolls |

Ein Fehler, der nur eine Verbindung betrifft, beendet diese Verbindung. Das
Werkzeug läuft weiter und liefert den Exit-Code des ersten aufgetretenen Fehlers
erst, wenn Sie es beenden.

## 6. Rollen und Rechte

Das Werkzeug kennt keine Benutzer, Rollen oder Anmeldung. Beim Aufzeichnen
leitet es die Anmeldung an die Datenbank weiter; beim Wiedergeben ist die
Anmeldung keine Sicherheitsgrenze, verlassen Sie sich nicht darauf, dass
Zugangsdaten geprüft werden. Betreiben Sie es nur in
einer kontrollierten Testumgebung, und lassen Sie es nur auf der Adresse
lauschen, die Sie mit `--listen` angegeben haben.

## 7. Fehlerbehebung

Jede Meldung des Werkzeugs trägt einen Code der Form `PGR-E…` (Fehler) oder
`PGR-W…` (Warnung). Der Text beginnt mit der Fehlerklasse und dem Code, zum
Beispiel `Replay [PGR-E5001]: …`.

### Fehlercodes

| Code | Bedeutung | Ursache und Lösung |
|---|---|---|
| `PGR-E1000` | sonstiger Fehler | Unerwarteter Fehler. Starten Sie mit `--log-level debug` neu, und melden Sie das Problem mit der Ausgabe. |
| `PGR-E2000`, `PGR-E2001` | ungültiger Aufruf | Eine Option fehlt, ist unbekannt, hat einen ungültigen Wert oder passt nicht zu einer anderen Option. Prüfen Sie den Aufruf mit `--help`. |
| `PGR-E2002` | Zieldatei existiert bereits | Wählen Sie einen anderen Dateinamen, oder ergänzen Sie `--force`, um die Datei zu ersetzen. |
| `PGR-E2003` | zeitgetreues Einspielen ohne Zeitangaben | Mindestens einer Anfrage der Aufzeichnung fehlt die Zeitangabe. Zeichnen Sie mit `--record-timing` erneut auf, oder starten Sie ohne `--keep-timing`. |
| `PGR-E2004` | Konfigurationsdatei nicht lesbar oder ungültig | Die Meldung nennt den Schlüssel oder die Verbindung. Prüfen Sie YAML, Schlüssel, Abschnitt, Werte und `sslmode` (erlaubt sind `disable` und `require`). |
| `PGR-E2005` | Umgebungsvariable eines Platzhalters nicht gesetzt | Setzen Sie die Variable, die als `${VAR}` in der benutzten Verbindung steht. |
| `PGR-E2006` | Klartext-Passwort in der Konfigurationsdatei | Ersetzen Sie das Passwort in der URL durch einen Platzhalter `${VAR}`. |
| `PGR-E2007` | Zertifikat, Schlüssel oder Zertifizierungsstelle nicht verwendbar | Die Datei fehlt, ist nicht lesbar oder kein gültiges PEM, oder Zertifikat und Schlüssel gehören nicht zusammen. Prüfen Sie `--tls-cert`, `--tls-key` und `--upstream-ca`; ein abgelaufenes eigenes Zertifikat und ein Schlüssel mit Passwort sind nicht zulässig. Läuft ein Zertifikat der Datenbank oder der Zertifizierungsstelle ab, meldet das Werkzeug beim Verbinden `PGR-E4005`. |
| `PGR-E3000`, `PGR-E3001` | Aufzeichnung nicht lesbar oder nicht schreibbar | Die Datei fehlt, oder Sie haben keine Rechte. Prüfen Sie Pfad und Dateirechte. |
| `PGR-E3002` | unbekannte Version der Aufzeichnung | Die Datei stammt aus einer anderen Programmversion. Zeichnen Sie mit der verwendeten Version erneut auf. |
| `PGR-E3003` | Aufzeichnung beschädigt | Die Datei ist unvollständig oder verändert. Zeichnen Sie erneut auf. |
| `PGR-E3004` | Aufzeichnung ohne verwendbare Sitzung | Beim Aufzeichnen hat keine Verbindung eine Anfrage gestellt. Zeichnen Sie erneut auf. |
| `PGR-E4000`, `PGR-E4003` | Verbindung unerwartet beendet | Die Verbindung brach mitten in einer Anfrage ab. Prüfen Sie Netzwerk, Datenbank und Anwendung. |
| `PGR-E4001` | Adresse nicht nutzbar | Der Port aus `--listen` ist belegt oder nicht erlaubt. Wählen Sie einen freien Port. |
| `PGR-E4002` | Datenbank nicht erreichbar | Prüfen Sie `--upstream`, die Datenbank und das Netzwerk. |
| `PGR-E4004` | Datenbank beantwortet eine eingespielte Anfrage mit einem Fehler (ohne `--compare-responses`) | Die Meldung nennt die Anfrage und die Antwort der Datenbank. Prüfen Sie Benutzer, Rechte und den Zustand der Datenbank, oder starten Sie mit `--continue-on-error`. |
| `PGR-E4005` | Anmeldung an der Datenbank fehlgeschlagen oder nicht unterstützt, Zertifikat der Datenbank oder der Zertifizierungsstelle ungültig oder abgelaufen, oder die Datenbank verlangt Verschlüsselung | Prüfen Sie Benutzer und Passwort (`PGWIRE_RECORDER_PASSWORD`). Unterstützt sind Klartext-Passwort, MD5 und SCRAM-SHA-256. Setzen Sie `--upstream-tls`, wenn die Datenbank Verschlüsselung verlangt. |
| `PGR-E5000`, `PGR-E5001` | Abweichung bei der Wiedergabe | Ihre Anwendung hat eine andere Anfrage gestellt als aufgezeichnet. Die Meldung nennt die erwartete und die empfangene Anfrage. Zeichnen Sie erneut auf, oder korrigieren Sie die Anwendung. |
| `PGR-E5002` | aufgezeichnete Anfragen oder Sitzungen nicht verbraucht | Ihr Test hat weniger Anfragen gestellt oder weniger Verbindungen geöffnet als aufgezeichnet, und `--fail-on-unconsumed` ist gesetzt. |
| `PGR-E5003` | Anfrage ohne aufgezeichnete Sitzung | Ihre Anwendung hat auf mehr Verbindungen Anfragen gestellt, als Sitzungen aufgezeichnet sind. Zeichnen Sie den Ablauf erneut auf, oder öffnen Sie weniger Verbindungen. |
| `PGR-E5004` | Antwort der Datenbank weicht von der Aufzeichnung ab | Nur mit `--compare-responses`. Die Meldung nennt Sitzung, Anfrage und die Art der Abweichung. Prüfen Sie, ob die Datenbank dieselbe Struktur liefert wie bei der Aufzeichnung (Tabellen, Spalten, Rechte). |
| `PGR-E6000`, `PGR-E6001` | nicht unterstützte Nachricht | Die Anwendung nutzt eine Funktion, die das Werkzeug nicht unterstützt, zum Beispiel `COPY`. Verwenden Sie diese Funktion im aufgezeichneten Ablauf nicht. |
| `PGR-E6002` | nicht unterstützte Protokollversion | Das Werkzeug unterstützt Version 3.0 des PostgreSQL-Protokolls. Verwenden Sie einen Treiber, der sie nutzt. |
| `PGR-E6003` | unverschlüsselte Verbindung nicht zugelassen | Das Werkzeug läuft mit `--tls-cert`, und die Anwendung hat ohne Verschlüsselung verbunden. Schalten Sie die Verschlüsselung in der Anwendung ein, oder starten Sie mit `--allow-plaintext`. |

### Warnungen

| Code | Bedeutung | Hinweis |
|---|---|---|
| `PGR-W2001` | Wiedergabe endete vor der letzten aufgezeichneten Anfrage | Ihr Test hat nicht alle aufgezeichneten Anfragen ausgeführt. Mit `--fail-on-unconsumed` wird das zum Fehler. |
| `PGR-W3001` | Abbruchwunsch nicht weitergeleitet | Die Anwendung hat versucht, eine laufende Anfrage abzubrechen. Das Werkzeug leitet diesen Wunsch nicht weiter und schließt die Verbindung. |
| `PGR-W3002` | Verschlüsselung einer Verbindung gescheitert | Die Anwendung hat die Aushandlung abgebrochen oder das Zertifikat nicht akzeptiert. Prüfen Sie, ob die Anwendung dem Zertifikat des Werkzeugs vertraut. |

### Die Anwendung kann sich nicht verbinden

Fordert Ihr Treiber eine verschlüsselte Verbindung an und bricht ab, wenn sie
abgelehnt wird, schalten Sie die Verschlüsselung in der Verbindungszeichenfolge
ab (zum Beispiel `sslmode=disable`), oder starten Sie das Werkzeug mit
`--tls-cert` und `--tls-key`. Vertraut Ihr Treiber dem Zertifikat nicht, hinterlegen
Sie dessen Zertifizierungsstelle im Treiber.

## 8. FAQ

**Brauche ich für die Wiedergabe eine Datenbank?**
Nein. Die Wiedergabe beantwortet alle aufgezeichneten Anfragen aus der Datei.

**Kann ich eine Aufzeichnung in die Versionsverwaltung legen?**
Ja. Die Standard-Datei ist Text (YAML) und ändert sich nur, wenn sich die
Aufzeichnung ändert; eine SQLite-Aufzeichnung ist binär. Beachten Sie den Hinweis zu vertraulichen Daten in der Einleitung.

**Kann ich die Datei von Hand bearbeiten?**
Das Werkzeug unterstützt das nicht. Zeichnen Sie stattdessen erneut auf.

**Funktioniert die Wiedergabe, wenn ich die Anfragen umformuliere?**
Nein. Die Anfrage muss Zeichen für Zeichen der aufgezeichneten entsprechen,
auch in Leerzeichen, Kommentaren und Groß- und Kleinschreibung. Beim erweiterten
Protokoll gilt das auch für Namen, Parametertypen, Parameterwerte, Formate und
das Zeilenlimit.

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

Es gibt noch keine veröffentlichte Version.
