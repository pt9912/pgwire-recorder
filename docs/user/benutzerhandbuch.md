# Benutzerhandbuch: pgwire-recorder

Version: 0.5  
Software-Version: `dev` (Ausgabe von `pgwire-recorder version`)  
Stand: 12.10.2026  
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
  die echte Datenbank weiter und speichert die Kommunikation in einer Datei, auf
  Wunsch als SQLite-Datei statt als YAML (siehe
  [Das Format der Aufzeichnung festlegen](#das-format-der-aufzeichnung-festlegen))
  und mit der Zeit seit der Startup-Nachricht, zu der jede Anfrage eintraf (siehe
  [Zeitangaben aufzeichnen und einspielen](#zeitangaben-aufzeichnen-und-einspielen)).
* **Wiedergeben (`replay`):** Ihre Anwendung verbindet sich wieder mit
  `pgwire-recorder`. Das Werkzeug beantwortet die Anfragen aus der Datei. Eine
  Datenbank ist nicht nötig.
* **Einspielen (`play`):** Das Werkzeug führt die aufgezeichneten Anfragen aus der
  Datei gegen eine Datenbank aus. Ihre Anwendung ist nicht beteiligt. Auf Wunsch wartet
  es auf die aufgezeichneten Abstände zwischen den Anfragen (siehe
  [Zeitangaben aufzeichnen und einspielen](#zeitangaben-aufzeichnen-und-einspielen)) und
  vergleicht die Antworten der Datenbank mit der Aufzeichnung (siehe
  [Antworten beim Einspielen vergleichen](#antworten-beim-einspielen-vergleichen)).

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

Dieses Beispiel schreibt YAML. Eine SQLite-Datei zeichnen Sie mit `--format sqlite`
auf (siehe [Das Format der Aufzeichnung festlegen](#das-format-der-aufzeichnung-festlegen)).

## 4. Aufgaben ausführen

### Eine Anwendung aufzeichnen

#### Voraussetzung

Die Datenbank ist erreichbar, und die Zieldatei existiert noch nicht.

#### Vorgehen

1. Starten Sie `pgwire-recorder record` mit den Optionen `--listen`,
   `--upstream` und `--output`; für eine SQLite-Datei statt YAML ergänzen Sie
   `--format sqlite`.
2. Verbinden Sie Ihre Anwendung mit der Adresse aus `--listen`.
3. Führen Sie den gewünschten Ablauf aus.
4. Beenden Sie das Werkzeug mit `Strg+C` oder `SIGTERM`.

#### Ergebnis

Die Datei aus `--output` enthält die aufgezeichnete Kommunikation, ohne `--format`
als YAML. Sie wird nach dem Ende jeder Verbindung und beim Beenden aktualisiert.

#### Hinweise

* Existiert die Zieldatei bereits, bricht das Werkzeug ab (`PGR-E2002`), gleich in
  welchem Format sie vorliegt und welches Format Sie wählen. Wollen
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
* Bis zur ersten Aktualisierung (nach der ersten beendeten Verbindung oder beim
  Beenden) bleibt eine vorhandene Zieldatei mit `--force` unverändert. Das Werkzeug
  schreibt die Aufzeichnung zuerst in
  eine temporäre Datei `.<Name der Zieldatei>.<16 Hexziffern>.tmp` im Verzeichnis
  der Zieldatei und ersetzt die Zieldatei danach in einem Schritt. Die Zieldatei
  ist damit vollständig oder unverändert, auch nach dem zwangsweisen Beenden von
  Verbindungen beim Herunterfahren. Das gilt für `yaml` bei jeder Aktualisierung; bei
  `sqlite` legt der erste Schreibvorgang die Datei auf diesem Weg an, und jede weitere
  beendete Verbindung wird danach in der Zieldatei selbst ergänzt (siehe
  [Das Format der Aufzeichnung festlegen](#das-format-der-aufzeichnung-festlegen)). Schlägt das Schreiben fehl, meldet das
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

### Das Format der Aufzeichnung festlegen

`record` schreibt die Aufzeichnung in einem von zwei Formaten: `yaml` (Standard) ist
eine Textdatei, `sqlite` ist eine SQLite-Datei, die Sie mit `sqlite3` oder einem
anderen SQLite-Werkzeug abfragen können. `replay` und `play` lesen beide.

#### Voraussetzung

Wie beim Aufzeichnen: Die Datenbank ist erreichbar, und die Zieldatei existiert noch
nicht. Das Verzeichnis der Zieldatei muss für das Werkzeug beschreibbar sein.

#### Vorgehen

1. Nennen Sie das Format beim Aufzeichnen:

   ```bash
   pgwire-recorder record \
     --listen 0.0.0.0:15432 \
     --upstream postgres:5432 \
     --output ./recordings/users.sqlite \
     --format sqlite
   ```

   Statt der Option genügt die Umgebungsvariable `PGWIRE_RECORDER_FORMAT=sqlite` oder
   der Schlüssel `format` im Abschnitt `record:` der Konfigurationsdatei:

   ```yaml
   record:
     listen: 0.0.0.0:15432
     upstream: postgres:5432
     output: ./recordings/users.sqlite
     format: sqlite
   ```

2. Zeichnen Sie wie gewohnt auf, und beenden Sie das Werkzeug.
3. Geben Sie die Datei wieder oder spielen Sie sie ein, ohne das Format zu nennen:

   ```bash
   pgwire-recorder replay --listen 0.0.0.0:15432 --input ./recordings/users.sqlite
   pgwire-recorder play --upstream postgres:5432 --input ./recordings/users.sqlite
   ```

#### Ergebnis

Die Datei aus `--output` ist eine SQLite-Datei. `replay` und `play` verhalten sich mit
ihr wie mit einer YAML-Aufzeichnung desselben Ablaufs: Ihre Anwendung erhält dieselben
Antworten, und `play` führt dieselben Anfragen aus.

#### Hinweise

* Die Werte von `--format` sind `yaml` und `sqlite`, in Kleinbuchstaben. Jeder andere
  Wert, auch `SQLite`, `sqlite3` und ein leerer Wert auf der Kommandozeile, ist ein
  ungültiger Aufruf (`PGR-E2001`); in der Konfigurationsdatei ist er `PGR-E2004`.
  Der Standard ist `yaml`. Option, Umgebungsvariable und Konfigurationsdatei folgen der
  Priorität aus [Einstellungen](#5-einstellungen); eine leere Umgebungsvariable gilt als
  nicht gesetzt, eine mit ungültigem Wert ist `PGR-E2001`, auch wenn Sie die Option
  zugleich angeben. `replay` und `play` kennen `--format` nicht (`PGR-E2001`).
* `replay` und `play` erkennen das Format an den ersten Bytes der Datei, nicht an
  ihrem Namen: Eine SQLite-Datei mit der Endung `.yaml` liest das Werkzeug als SQLite,
  eine YAML-Datei mit der Endung `.sqlite` als YAML. Eine Datei, die nicht mit der
  Kennung einer SQLite-Datei beginnt, liest es als YAML und meldet `PGR-E3003`, wenn
  sie als YAML nicht lesbar ist; das gilt auch für eine leere Datei.
* Eine vorhandene Zieldatei ersetzt `--force` vollständig; das Werkzeug ergänzt sie nie.
  Eine YAML-Datei ersetzt es mit `--format sqlite` durch eine SQLite-Datei und umgekehrt.
* Bei `sqlite` entsteht die Datei mit dem ersten Schreibvorgang, nach der ersten
  beendeten Verbindung oder beim Beenden; vorher gibt es den Zielpfad nicht. Jede weitere
  beendete Verbindung ergänzt das Werkzeug danach in der Datei. Zeichnen Sie ohne jede
  Verbindung auf, entsteht eine gültige Datei ohne Sitzung; `replay` meldet sie beim
  Start als `PGR-E3004`.
* Die Datei lässt sich während des Aufzeichnens lesen: `sqlite3` und `replay` sehen die
  Verbindungen, die zu dem Zeitpunkt beendet waren.
* `replay` und `play` öffnen die Datei nur lesend und verändern sie nicht, auch auf einem
  schreibgeschützten Verzeichnis.

#### Randfälle

* Liegt beim Start von `record` eine Datei `<Zieldatei>-journal` neben dem Zielpfad,
  endet `record` mit `PGR-E3001`, auch mit `--force` und auch ohne Zieldatei; mit der
  Zieldatei und ohne `--force` gilt zuerst `PGR-E2002`. Nennt die Meldung die Datei ein
  Rollback-Journal eines früheren Laufs und hat sie Inhalt: Öffnen Sie die Zieldatei einmal mit
  `sqlite3` (eine Abfrage genügt), dann macht SQLite den Abbruch rückgängig und entfernt
  die Journal-Datei; danach startet `record`. Eine leere Journal-Datei entfernt
  `sqlite3` nicht; löschen Sie sie selbst, danach startet `record`.
* Hält ein anderer Prozess die Datei gesperrt, wartet das Lesen durch `replay` und
  `play` höchstens fünf Sekunden und endet dann mit `PGR-E3001`.
* Eine SQLite-Datei ist binär. Ist sie keine gültige Aufzeichnung dieses Werkzeugs (eine
  fremde SQLite-Datei, eine fehlende Tabelle oder Spalte, eine abgeschnittene Datei,
  ein widersprüchlicher Inhalt), meldet das Werkzeug `PGR-E3003`; eine unbekannte
  Version der Datei ist `PGR-E3002`.
* Das Aufzeichnen und Lesen von SQLite-Dateien ist unter Linux im Container auf einem
  lokalen Dateisystem geprüft. Für andere Systeme und für Netzwerk-Dateisysteme ist es
  nicht geprüft.

#### Eine SQLite-Aufzeichnung abfragen

Die Datei hat vier Tabellen:

| Tabelle | Inhalt |
|---|---|
| `meta` | Zeilen mit `key` und `value`: `format` (`pgwire-recorder`) und `version` |
| `session` | eine Zeile je aufgezeichneter Verbindung: `id` (ab 1), `startup` und `server_parameters` als JSON-Text |
| `interaction` | eine Zeile je Anfrage oder Folge des erweiterten Protokolls: `session_id`, `sequence`, `type` (`query` oder `extended`), `offset_ms` (die Zeitangabe, ohne `--record-timing` ohne Wert, also `NULL`) und `sql` (der Text einer einfachen Anfrage, bei `extended` ohne Wert, also `NULL`) |
| `message` | die Nachrichten einer Interaktion: `session_id`, `sequence`, `position`, `direction` (`client` oder `server`), `group_no`, `kind` und `fields` als JSON-Text |

Bei einer einfachen Anfrage stehen in `message` nur die Antworten des Servers; die
Anfrage steht in `interaction.sql`. Der SQL-Text einer erweiterten Anfrage steht im Feld
`sql` der Nachricht mit `kind = 'parse'`. Zum Beispiel:

```bash
sqlite3 ./recordings/users.sqlite "select session_id, sequence, sql from interaction where type = 'query'"
sqlite3 ./recordings/users.sqlite "select session_id, sequence, json_extract(fields, '$.sql') from message where kind = 'parse'"
```

Der zweite Aufruf braucht ein `sqlite3` mit JSON-Funktionen. Ändern Sie die Datei nicht;
das Werkzeug unterstützt das nicht.

### Eine Aufzeichnung wiedergeben

#### Voraussetzung

Eine Aufzeichnung mit einer Verbindung liegt vor, als YAML- oder als SQLite-Datei. Das
Werkzeug erkennt das Format selbst (siehe
[Das Format der Aufzeichnung festlegen](#das-format-der-aufzeichnung-festlegen)).

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
* Trägt die Aufzeichnung Zeitangaben (`offset_ms`), wartet `replay` nie auf sie; es
  prüft sie beim Laden (siehe
  [Zeitangaben aufzeichnen und einspielen](#zeitangaben-aufzeichnen-und-einspielen)).

### Eine Aufzeichnung in eine Datenbank einspielen

Damit führen Sie die aufgezeichneten Anfragen erneut aus, zum Beispiel um eine
Komponente zu testen, die Änderungen der Datenbank verarbeitet (Change Data Capture).

#### Voraussetzung

Eine Aufzeichnung liegt vor, als YAML- oder als SQLite-Datei, und die Zieldatenbank ist
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

* Die drei folgenden Hinweise zu Fehlerantworten der Datenbank beschreiben das
  Einspielen ohne `--compare-responses`; mit der Option gilt
  [Antworten beim Einspielen vergleichen](#antworten-beim-einspielen-vergleichen).
  Ebenso beschreiben die Hinweise zu Anmeldung und TLS, zu COPY und zu
  `--finish-session-on-interrupt` das Einspielen ohne diese Option.
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
  Anfrage zählt nicht als Fehler. Der Exit-Code ist 0, nach einem früheren Fehler 4, nach
  einer Abweichung beim Vergleich 5.
* Das Einspielen führt Anfragen aus. Verwenden Sie es nicht gegen eine Datenbank,
  deren Inhalt Sie nicht verändern dürfen.

### Antworten beim Einspielen vergleichen

Mit `--compare-responses` vergleicht `play` nach jeder Interaktion die Antworten der
Datenbank mit den aufgezeichneten. Eine Interaktion ist eine einfache Anfrage oder
eine Folge des erweiterten Protokolls bis zu ihrem `Sync`; ihre Nummer ist die der
Aufzeichnung. So erkennen Sie, ob die Datenbank, in die Sie einspielen, anders
antwortet als die, gegen die Sie aufgezeichnet haben, etwa mit anderen Spalten.

#### Voraussetzung

Wie beim Einspielen.

#### Vorgehen

1. Starten Sie das Einspielen mit der Option:

   ```bash
   pgwire-recorder play \
     --upstream postgres:5432 \
     --input ./recordings/users.yaml \
     --compare-responses
   ```

   Statt der Option genügt die Umgebungsvariable
   `PGWIRE_RECORDER_COMPARE_RESPONSES=true` oder der Schlüssel `compare_responses: true`
   im Abschnitt `play:` der Konfigurationsdatei:

   ```yaml
   play:
     upstream: postgres:5432
     input: ./recordings/users.yaml
     compare_responses: true
   ```
2. Warten Sie, bis das Werkzeug endet.

#### Ergebnis

Stimmen die Antworten der Datenbank mit der Aufzeichnung überein, endet `play` mit
Exit-Code 0. Die Zeile zum Ende nennt, wie viele Interaktionen eingespielt und
verglichen wurden und wie viele abwichen:

```text
time=2026-10-12T14:03:12.481+02:00 level=INFO msg="play beendet" eingespielt=3 verglichen=3 abweichend=0
```

Weicht eine Antwort ab, meldet `play` `PGR-E5004` und endet mit Exit-Code 5:

```text
time=2026-10-12T14:03:12.482+02:00 level=ERROR msg=Fehler code=PGR-E5004 error="Replay [PGR-E5004]: Session 1, Interaktion 1: Antwort weicht von der Aufzeichnung ab: Position 1, Art Spaltenbeschreibung"
time=2026-10-12T14:03:12.482+02:00 level=INFO msg="play beendet" eingespielt=1 verglichen=1 abweichend=1
```

#### Hinweise

* Verglichen werden die Art und die Reihenfolge der Antworten der Datenbank. Bei der
  Beschreibung der Spalten zählen ihre Anzahl und je Spalte Name und Typ; beim
  erweiterten Protokoll zählt auch die Beschreibung der Parameter. Bei einem
  abgeschlossenen Befehl zählt der Befehl ohne die abschließenden Zahlen: `UPDATE 1` und
  `UPDATE 7` sind gleich, `UPDATE 1` und `DELETE 1` nicht. Bei einem Fehler zählt der
  SQLSTATE, bei der Bereitschaft für die nächste Anfrage der Transaktionsstatus. Das gilt
  für einfache Anfragen und für das erweiterte Protokoll.
* Nicht verglichen werden die Werte und die Zahl der Zeilen, Hinweise der Datenbank
  (`NOTICE`), Meldungen über Einstellungen des Servers (etwa nach
  `SET application_name`) und der Text eines Fehlers.
* Die Meldung `PGR-E5004` nennt die Sitzung (`Session`), die Nummer der Interaktion, die
  Position der ersten abweichenden Antwort und die Art der Abweichung. Die Position zählt
  ab 1; Zeilen, Hinweise und Meldungen über Einstellungen zählen nicht mit. Die Art ist
  `Nachrichtenart` (die Datenbank antwortet mit einer anderen Art von Antwort),
  `Spaltenbeschreibung`, `Befehl`, `Fehler` oder `Transaktionsstatus`. Bei der Art
  `Fehler` nennt die Meldung die SQLSTATE der Aufzeichnung und der Datenbank, nie den
  Text des Fehlers. Je Interaktion meldet `play` die erste Abweichung.
* Ohne `--continue-on-error` endet `play` nach der Interaktion mit der Abweichung; die
  weiteren Interaktionen zählen nicht als eingespielt. Mit `--continue-on-error` macht es
  mit der nächsten Interaktion weiter und meldet jede Interaktion mit einer Abweichung;
  der Exit-Code ist dann 5. Ohne `--compare-responses` findet kein Vergleich statt.
* Antwortet die Datenbank mit einem Fehler, den die Aufzeichnung an dieser Stelle nicht
  trägt, ist das mit Vergleich `PGR-E5004` (Art `Nachrichtenart`) und nicht `PGR-E4004`.
  Trägt die Aufzeichnung an dieser Stelle eine Fehlerantwort mit demselben SQLSTATE, ist
  das keine Abweichung, auch bei einem anderen Text; `play` macht mit der nächsten
  Interaktion weiter, ohne `--continue-on-error`. Ein anderer SQLSTATE ist `PGR-E5004`
  mit der Art `Fehler`.
* Mit Vergleich hat `--allow-recorded-errors` keine Wirkung: Ein Fehler ohne Vorbild oder
  mit einem anderen SQLSTATE bleibt `PGR-E5004`, einer mit demselben SQLSTATE ist auch
  ohne die Option erwartet.
* Beendet die Datenbank die Verbindung mit einer Fehlerantwort `FATAL` oder `PANIC`, ist
  das mit Vergleich `PGR-E5004`, auch wenn die Aufzeichnung an dieser Stelle einen Fehler
  mit demselben SQLSTATE trägt. Die Art ist `Fehler` an der Position der Fehlerantwort;
  weicht davor schon eine Antwort ab, nennt die Meldung diese Abweichung. Ohne Vergleich
  ist es `PGR-E4003`.
* Bricht die Verbindung ohne Fehlerantwort der Datenbank ab, ist das `PGR-E4003` mit
  Exit-Code 4, auch nach einer Abweichung; beide Meldungen stehen im Log.
* Das erste Abbruchsignal (`Strg+C`, `SIGTERM`) beendet die laufende Interaktion mit
  Vergleich; eine Abweichung darin ist `PGR-E5004`, und der Exit-Code ist 5, wenn sie oder
  eine frühere Interaktion abwich.
* Die Zeile `play beendet` trägt die Zähler `eingespielt`, `verglichen` und `abweichend`.
  Eine Interaktion, die ein Verbindungsverlust unterbricht, zählt als eingespielt, nicht
  als verglichen. Ohne `--compare-responses` sind `verglichen` und `abweichend` 0. Die
  Zeile hat die Stufe `info`; mit `--log-level warn` erscheint sie nicht.
* Die Option gibt es nur bei `play`: bei `record` und `replay` ist sie ein ungültiger
  Aufruf (`PGR-E2001`), ebenso `PGWIRE_RECORDER_COMPARE_RESPONSES`. Der Schlüssel
  `compare_responses` im Abschnitt `record:` oder `replay:` ist `PGR-E2004`, ebenso ein
  anderer Wert als `true` oder `false`. Die Umgebungsvariable mit einem anderen Wert ist
  `PGR-E2001`.

### Zeitangaben aufzeichnen und einspielen

Mit `--record-timing` hält `record` an jeder Anfrage fest, wie viele Millisekunden nach
der Startup-Nachricht der Verbindung sie eintraf. Mit `--keep-timing` wartet `play` vor den
Anfragen auf diese Zeit; beim Standard `relative` sendet es die erste Anfrage jeder
Verbindung ohne Wartezeit. Ohne `--keep-timing` wartet `play` nie, auch wenn die
Aufzeichnung Zeitangaben trägt.

#### Voraussetzung

Wie beim Aufzeichnen und beim Einspielen. Für `--keep-timing` muss jede Anfrage der
Aufzeichnung eine Zeitangabe tragen; zeichnen Sie dafür mit `--record-timing` auf.

#### Vorgehen

1. Zeichnen Sie mit `--record-timing` auf:

   ```bash
   pgwire-recorder record \
     --listen 0.0.0.0:15432 \
     --upstream postgres:5432 \
     --output ./recordings/users.yaml \
     --record-timing
   ```

   Statt der Option genügt die Umgebungsvariable `PGWIRE_RECORDER_RECORD_TIMING=true`
   oder der Schlüssel `record_timing: true` im Abschnitt `record:` der
   Konfigurationsdatei.
2. Zeichnen Sie wie gewohnt auf, und beenden Sie das Werkzeug.
3. Spielen Sie die Aufzeichnung mit ihren Zeitangaben ein:

   ```bash
   pgwire-recorder play \
     --upstream postgres:5432 \
     --input ./recordings/users.yaml \
     --keep-timing
   ```

   Statt der Option genügt `PGWIRE_RECORDER_KEEP_TIMING=true` oder der Schlüssel
   `keep_timing: true` im Abschnitt `play:`. Diese Konfigurationsdatei setzt alle
   vier Einstellungen der Zeitangaben; `timing_mode` und `timing_reference` erklären die
   Hinweise zum Einspielen:

   ```yaml
   record:
     record_timing: true
   play:
     keep_timing: true
     timing_mode: absolute
     timing_reference: first-request
   ```

#### Ergebnis

Jede aufgezeichnete Anfrage trägt das Feld `offset_ms`. In YAML steht es neben
`sequence`:

```yaml
      - sequence: 2
        offset_ms: 1012
        request:
```

In SQLite steht es in der Spalte `offset_ms` der Tabelle `interaction`:

```bash
sqlite3 ./recordings/users.sqlite "select session_id, sequence, offset_ms from interaction"
```

Beim Einspielen mit `--keep-timing` wartet `play` vor den Anfragen auf ihre aufgezeichnete
Zeit; beim Standard `relative` sendet es die erste Anfrage jeder Verbindung ohne Wartezeit.
Liegen die erste und die letzte Anfrage einer Aufzeichnung mit einer Verbindung eine Sekunde
auseinander, liegt zwischen den Log-Zeilen `play gestartet` und `play beendet` mit dem
Standard `relative` mindestens eine Sekunde; ohne `--keep-timing` liegt weniger dazwischen.

#### Hinweise zum Aufzeichnen

* `offset_ms` ist eine ganze Zahl ab 0: die Millisekunden seit der Startup-Nachricht der
  Verbindung bis zum Eintreffen der Anfrage.
* Die Zeit, die Ihre Anwendung bis zum Senden der Startup-Nachricht braucht, zählt nicht
  mit, auch nicht nach einer Anfrage nach Verschlüsselung. Die Zeit zwischen der Anmeldung
  und der ersten Anfrage zählt mit.
* Der Wert der ersten Anfrage einer Verbindung ist nicht notwendig 0.
* Jede Verbindung zählt ab ihrem eigenen Beginn.
* Das Feld tragen die Anfragen beider Protokolle, einfach und erweitert, und eine
  Lebendprüfung (siehe [Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)).
* Ohne `--record-timing`, dem Standard, tragen die Anfragen kein Feld: In YAML fehlt
  `offset_ms`, in SQLite ist die Spalte ohne Wert, also `NULL`.
* Zwei Aufzeichnungen desselben Ablaufs können sich mit `--record-timing` in den Werten von
  `offset_ms` unterscheiden.
* `--record-timing` nimmt keinen Wert oder `=true` beziehungsweise `=false`; die
  Umgebungsvariable nimmt `true` oder `false`, leer gilt sie als nicht gesetzt. Jeder
  andere Wert ist ungültig: als Option, auch `=1`, `=True` und `=` ohne Wert, und als
  Umgebungsvariable `PGR-E2001`, als Schlüssel der Konfigurationsdatei `PGR-E2004`. Die
  Umgebungsvariable mit ungültigem Wert ist auch dann `PGR-E2001`, wenn Sie die Option
  zugleich angeben. Die Option geht der Umgebungsvariablen vor, diese der Datei.

#### Hinweise zum Einspielen

* `--keep-timing` nimmt keinen Wert oder `=true` beziehungsweise `=false`, die
  Umgebungsvariable `PGWIRE_RECORDER_KEEP_TIMING` `true` oder `false`; Fehlerfälle,
  leere Umgebungsvariable und Priorität wie bei `--record-timing`. Die Zeitangaben gelten
  für YAML- und SQLite-Aufzeichnungen.
* `--timing-mode` wählt, worauf `play` wartet; der Standard ist `relative`:
  * `relative`: Vor jeder Anfrage außer der ersten einer Verbindung wartet `play`, bis seit
    dem Beginn der vorigen Anfrage, dem Senden an die Datenbank, mindestens der Unterschied
    der aufgezeichneten Zeiten beider Anfragen vergangen ist. Dauert die vorige Anfrage
    länger als aufgezeichnet, sendet `play` die nächste, sobald sie beantwortet ist; die
    Verspätung holt es bei den folgenden Anfragen nicht auf.
  * `absolute`: Vor jeder Anfrage wartet `play`, bis seit einem Bezugspunkt mindestens die
    aufgezeichnete Zeit der Anfrage vergangen ist. Ist `play` verspätet, holt es auf: Der
    Abstand zwischen zwei Anfragen kann kürzer sein als aufgezeichnet.
* `--timing-reference` wählt den Bezugspunkt von `absolute`; der Standard ist `connect`:
  `connect` ist der Beginn des Verbindungsaufbaus zur Datenbank, `first-request` der
  Beginn der ersten Anfrage der Verbindung. Mit `connect` wartet `play` auch vor der ersten
  Anfrage, bis deren aufgezeichnete Zeit erreicht ist. Mit `first-request` und mit
  `relative` sendet es die erste Anfrage jeder Verbindung ohne Wartezeit.
* Jede Verbindung der Aufzeichnung beginnt mit einer eigenen Uhr; `play` spielt die
  Verbindungen nacheinander ein.
* Die Werte sind `relative` und `absolute` sowie `connect` und `first-request`, in
  Kleinbuchstaben. Jeder andere Wert ist ein ungültiger Aufruf (`PGR-E2001`), in der
  Konfigurationsdatei `PGR-E2004`, auch `Absolute`, `first_request` und ein leerer Wert;
  die Meldung nennt den Wert der Umgebungsvariablen und der Datei nicht. Eine leere
  Umgebungsvariable gilt als nicht gesetzt.

#### Ungültige Kombinationen

Jeder dieser Aufrufe endet mit `PGR-E2001` (Exit-Code 2):

* `--timing-mode` oder `--timing-reference` ohne `--keep-timing`.
* `--timing-reference` zusammen mit `relative`.
* Als gesetzt gilt, was Option, Umgebungsvariable oder Datei nennt, auch den Wert des
  Standards: `--timing-reference connect` ist mit `relative`, auch mit dem Standard
  `relative`, ungültig, und `--keep-timing=false` ist mit `--timing-mode absolute` ungültig.
* `record` kennt `--keep-timing`, `--timing-mode` und `--timing-reference` nicht, `play`
  kennt `--record-timing` nicht, und `replay` kennt keine der vier Optionen.

Die Reihenfolge der Optionen ändert nichts.

#### Fehlende Zeitangaben und Abbruch

* Trägt eine Anfrage der Aufzeichnung kein `offset_ms`, endet `play --keep-timing` beim Start
  mit `PGR-E2003` (Exit-Code 2), bevor es eine Verbindung öffnet. Die Meldung nennt die
  Sitzung und die Nummer der ersten Anfrage ohne Zeitangabe, in der Reihenfolge der Sitzungen
  und Anfragen, keinen Wert. Ohne `--keep-timing` spielt `play` eine solche Aufzeichnung
  ein.
* Das erste Abbruchsignal (`Strg+C` oder `SIGTERM`) beendet jede Wartezeit sofort. Ohne
  `--finish-session-on-interrupt` sendet `play` die folgenden Anfragen nicht und endet; mit
  der Option spielt es die übrigen Anfragen der Sitzung ohne Wartezeit ein und endet dann.
  Der Exit-Code ist 0, nach einem früheren Fehler 4, nach einer Abweichung beim Vergleich 5.

#### Prüfung beim Laden

`replay` und `play` prüfen die Zeitangaben jeder Aufzeichnung beim Laden, mit und ohne
`--keep-timing`. Eine Aufzeichnung ist beschädigt (`PGR-E3003`, Exit-Code 3), wenn

* ein `offset_ms` keine ganze Zahl ab 0 ist. In YAML sind `0` und Ziffern ohne führende Null
  bis 9223372036854775807 gültig; jede andere Schreibweise ist beschädigt, etwa `1.5`, `5.0`,
  `1e3`, `0x10`, `007`, `+5`, `-1`, `"5"`, `true`, `null` und ein Wert über
  9223372036854775807. In SQLite muss der Wert `NULL` oder eine ganze Zahl ab 0 sein; eine
  negative Zahl, ein Text und eine Dezimalzahl sind beschädigt.
* ein Wert kleiner ist als der vorige Wert derselben Sitzung.

Gleiche Werte hintereinander sind gültig, ebenso ein kleinerer Wert am Anfang der nächsten
Sitzung. Eine Anfrage ohne `offset_ms` zwischen zwei Anfragen mit Wert ist beim Laden gültig;
`play --keep-timing` lehnt sie ab (`PGR-E2003`). Die Meldung nennt die Nummer der Sitzung und
der Anfrage, nie einen Wert aus der Datei. Ist eine Anfrage beschädigt, meldet das Werkzeug
`PGR-E3003`, auch wenn einer anderen Anfrage das Feld fehlt und `play --keep-timing`
sonst `PGR-E2003` meldete.

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

Die Umgebungsvariable einer Option, die die Betriebsart nicht kennt (Spalte
*Betriebsart* der Tabelle), ist ein ungültiger Aufruf (`PGR-E2001`), sobald sie
gesetzt und nicht leer ist, gleich welchen Wert sie hat, ob `true`, `5s` oder ein Wert,
den die Option nicht annimmt; die Meldung nennt den Namen der Variable, nie ihren
Wert. So ist `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` bei `record` und `play` ein
Fehler, `PGWIRE_RECORDER_LISTEN` bei `play` und `PGWIRE_RECORDER_COMPARE_RESPONSES` bei
`record` und `replay`. Wer `record` und `play` in derselben
Umgebung ausführt, setzt dort die Variablen der Optionen der jeweils anderen
Betriebsart nicht. `config show` prüft von den Variablen nur
`PGWIRE_RECORDER_CONFIG`. Unbeachtet bleiben: `PGWIRE_RECORDER_PASSWORD`, das keine Option
hat, bei `record` und `replay`, ferner bei `record`, `replay` und `play` eine Variable mit
dem Präfix `PGWIRE_RECORDER_`, deren Name zu keiner Option gehört und die nicht
`PGWIRE_RECORDER_PASSWORD` heißt, und
`pgwire_recorder_listen`, eine andere Schreibweise des Namens einer Option.

| Option | Betriebsart | Umgebungsvariable | Standard |
|---|---|---|---|
| `--listen` | `record`, `replay` | `PGWIRE_RECORDER_LISTEN` | Pflicht |
| `--upstream` | `record`, `play` | `PGWIRE_RECORDER_UPSTREAM` | Pflicht (`host:port` oder Name einer Verbindung) |
| `--output` | `record` | `PGWIRE_RECORDER_OUTPUT` | Pflicht |
| `--force` | `record` | `PGWIRE_RECORDER_FORCE` | `false` |
| `--format` | `record` | `PGWIRE_RECORDER_FORMAT` | `yaml` (Werte `yaml`, `sqlite`) |
| `--record-timing` | `record` | `PGWIRE_RECORDER_RECORD_TIMING` | `false` |
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
| `--keep-timing` | `play` | `PGWIRE_RECORDER_KEEP_TIMING` | `false` |
| `--timing-mode` | `play` | `PGWIRE_RECORDER_TIMING_MODE` | `relative` (Werte `relative`, `absolute`; nur mit `--keep-timing`) |
| `--timing-reference` | `play` | `PGWIRE_RECORDER_TIMING_REFERENCE` | `connect` (Werte `connect`, `first-request`; nur mit `--keep-timing` und `absolute`) |
| `--compare-responses` | `play` | `PGWIRE_RECORDER_COMPARE_RESPONSES` | `false` |
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
| 5 | Abweichung bei der Wiedergabe, mit `--fail-on-unconsumed` auch nicht gestellte Anfragen oder nie verwendete Sitzungen, bei `play` mit `--compare-responses` eine Antwort der Datenbank, die von der Aufzeichnung abweicht |
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
| `PGR-E2001` | ungültiger Aufruf | Eine Option fehlt, ist unbekannt oder hat einen ungültigen Wert (bei `--format` jeder außer `yaml` und `sqlite`), oder eine gesetzte, nicht leere Umgebungsvariable gehört zu einer Option, die die Betriebsart nicht kennt (die Meldung nennt den Namen der Variable, nicht den Wert; siehe [Einstellungen](#5-einstellungen)), oder `--upstream-ca` steht, obwohl `play` kein TLS nutzt (ohne `--upstream-tls` und ohne `sslmode=require` der Verbindung, oder mit `--upstream-tls=false`), oder `--timing-mode` oder `--timing-reference` steht ohne `--keep-timing`, oder `--timing-reference` steht mit `relative` (siehe [Zeitangaben aufzeichnen und einspielen](#zeitangaben-aufzeichnen-und-einspielen)). Prüfen Sie den Aufruf mit `--help`. |
| `PGR-E2002` | Zieldatei existiert bereits | Wählen Sie einen anderen Dateinamen, oder ergänzen Sie `--force`, um die Datei zu ersetzen. |
| `PGR-E2003` | Aufzeichnung ohne Zeitangabe bei `--keep-timing` | `play --keep-timing` hat eine Aufzeichnung, in der eine Anfrage kein `offset_ms` trägt. Die Meldung nennt die Sitzung und die Nummer der ersten Anfrage ohne Zeitangabe. Zeichnen Sie mit `--record-timing` erneut auf, oder spielen Sie ohne `--keep-timing` ein. |
| `PGR-E2004` | Konfigurationsdatei nicht lesbar oder ungültig | Die Meldung nennt den Schlüssel oder die Verbindung. Prüfen Sie YAML, Schlüssel, Abschnitt, Werte, `sslmode` (erlaubt sind `disable` und `require`), `format` (erlaubt sind `yaml` und `sqlite`), `record_timing`, `keep_timing` und `compare_responses` (erlaubt sind `true` und `false`), `timing_mode` (erlaubt sind `relative` und `absolute`) und `timing_reference` (erlaubt sind `connect` und `first-request`). |
| `PGR-E2005` | Umgebungsvariable eines Platzhalters nicht gesetzt | Die Meldung nennt die Verbindung und die erste fehlende Variable. Setzen Sie sie mit einem nicht leeren Wert; eine leere Variable gilt als nicht gesetzt. Bei `record` zählen nur die Variablen in Host und Port. |
| `PGR-E2006` | Klartext-Passwort in der Konfigurationsdatei | Die Meldung nennt die Verbindung. Ersetzen Sie das Passwort in der URL durch genau einen Platzhalter `${VAR}`, und entfernen Sie einen Parameter `password`. Das gilt für jede Verbindung der Datei, auch eine, die Sie nicht benutzen. |
| `PGR-E2007` | Datei von `--upstream-ca` nicht verwendbar | Die Datei fehlt, ist nicht lesbar, ist keine reguläre Datei oder enthält keine Zertifikate im PEM-Format. Die Meldung nennt die Option und den Grund, nicht den Pfad; der Exit-Code ist 2. Prüfen Sie Pfad, Rechte und Inhalt der Datei. |
| `PGR-E3001` | Aufzeichnung nicht lesbar oder nicht schreibbar | Die Datei fehlt, Sie haben keine Rechte, `--output` ist ein Verzeichnis, oder das Verzeichnis von `--output` fehlt. Bei `sqlite` außerdem: Neben dem Zielpfad liegt eine Datei `<Zielpfad>-journal`, oder eine SQLite-Datei ist länger als fünf Sekunden gesperrt. Prüfen Sie Pfad und Dateirechte, und siehe [Das Format der Aufzeichnung festlegen](#das-format-der-aufzeichnung-festlegen). |
| `PGR-E3002` | unbekannte Version der Aufzeichnung | Die Datei stammt aus einer anderen Programmversion; das gilt für beide Formate. Zeichnen Sie mit der verwendeten Version erneut auf. |
| `PGR-E3003` | Aufzeichnung beschädigt | Die Datei ist unvollständig oder verändert. Bei einer SQLite-Datei auch: Sie ist keine SQLite-Datei oder eine fremde, eine Tabelle oder Spalte fehlt, sie ist abgeschnitten, oder ihr Inhalt ist ungültig (zum Beispiel falsche Formatkennung, keine lückenlose Nummerierung, ein Eintrag ohne zugehörige Sitzung). Auch: Ein `offset_ms` ist keine ganze Zahl ab 0, oder ein Wert ist kleiner als der vorige derselben Sitzung; die Meldung nennt die Sitzung und die Anfrage, nie den Wert (siehe [Zeitangaben aufzeichnen und einspielen](#zeitangaben-aufzeichnen-und-einspielen)). Zeichnen Sie erneut auf. |
| `PGR-E3004` | Aufzeichnung ohne verwendbare Sitzung | `replay` meldet den Fehler beim Start. Beim Aufzeichnen hat keine Verbindung eine Anfrage gestellt, oder die Verbindungen haben nur Lebendprüfungen gesendet (Anfragen nur aus Leerraum und Kommentaren, siehe [Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)). Zeichnen Sie einen Ablauf mit mindestens einer anderen Anfrage auf. |
| `PGR-E4000` | Verbindung nicht anzunehmen | Das Werkzeug konnte eine eingehende Verbindung nicht annehmen, zum Beispiel weil zu viele Dateien offen sind. Es nimmt danach weiter Verbindungen an. Prüfen Sie die Grenzen des Systems. |
| `PGR-E4003` | Verbindung unerwartet beendet | Die Verbindung brach mitten in einer Anfrage ab. Prüfen Sie Netzwerk, Datenbank und Anwendung. |
| `PGR-E4001` | Adresse nicht nutzbar | Der Port aus `--listen` ist belegt oder nicht erlaubt. Wählen Sie einen freien Port. |
| `PGR-E4002` | Datenbank nicht erreichbar | Prüfen Sie `--upstream`, die Datenbank und das Netzwerk. Bei `play` mit TLS gehört dazu: Die Antwort auf die Anfrage nach Verschlüsselung ist weder `S` noch `N`, oder die Verbindung endet davor. |
| `PGR-E4004` | Datenbank beantwortet eine eingespielte Anfrage mit einem Fehler | Die Meldung nennt die Sitzung und die Nummer der Anfrage in der Aufzeichnung, dazu SQLSTATE und Meldung der Datenbank, nicht den Text der Anfrage. Prüfen Sie Benutzer, Rechte und den Zustand der Datenbank, oder starten Sie mit `--continue-on-error`. Mit `--compare-responses` ist ein Fehler, den die Aufzeichnung an dieser Stelle nicht trägt, `PGR-E5004`. |
| `PGR-E4005` | Anmeldung oder TLS zur Datenbank abgelehnt oder nicht unterstützt | Beim Einspielen hat die Datenbank die Anmeldung abgelehnt, etwa weil der Benutzer fehlt oder das Passwort falsch ist, oder sie verlangt ein Passwort, das `play` nicht hat, oder ein Anmeldeverfahren, das `play` nicht kann. Die Meldung nennt die Sitzung; bei einer Ablehnung außerdem SQLSTATE und Meldung der Datenbank, bei fehlendem Passwort, dass `play` keines hat. `play` meldet sich mit Klartext-Passwort, MD5 und SCRAM-SHA-256 an; bietet die Datenbank nur `SCRAM-SHA-256-PLUS` an oder verlangt sie Kerberos, GSSAPI oder SSPI, endet `play` ebenfalls mit diesem Code. Dasselbe gilt für eine Datenbank, die Verbindungen ohne Verschlüsselung ablehnt, wenn `play` kein TLS wünscht. Mit TLS endet `play` mit diesem Code auch, wenn die Datenbank TLS ablehnt (Antwort `N`, ohne Rückfall auf eine Verbindung ohne TLS) oder ihr Zertifikat nicht angenommen wird: Das Zertifikat lässt sich weder mit dem Speicher des Systems noch mit `--upstream-ca` überprüfen, der Host steht nicht im Zertifikat, oder das Zertifikat ist abgelaufen. Die Meldung nennt die Sitzung und den Grund, nie Pfad oder Inhalt der Datei aus `--upstream-ca`. Prüfen Sie Benutzer und Datenbank (`--user`, `--database`), das Passwort (siehe [Einstellungen](#5-einstellungen)), den Host der Verbindung und `--upstream-ca`. |
| `PGR-E4006` | Anfrage beim Beenden unvollständig | Die Frist `--shutdown-timeout` ist abgelaufen, oder ein zweites Signal hat sie ablaufen lassen, während eine Anfrage lief. Die Meldung nennt die Sitzung und die Anfrage; beim Aufzeichnen fehlt diese Anfrage in der Aufzeichnung. Lassen Sie die Anwendung vor dem Beenden zur Ruhe kommen, oder wählen Sie eine längere Frist (siehe [Herunterfahren mit Frist](#herunterfahren-mit-frist)). |
| `PGR-E5001` | Abweichung bei der Wiedergabe | Ihre Anwendung hat eine andere Anfrage gestellt als aufgezeichnet. Die Meldung nennt die erwartete und die empfangene Anfrage. Zeichnen Sie erneut auf, oder korrigieren Sie die Anwendung. |
| `PGR-E5002` | aufgezeichnete Anfragen oder Sitzungen nicht verbraucht | Ihr Test hat weniger Anfragen gestellt oder weniger Verbindungen geöffnet als aufgezeichnet, und `--fail-on-unconsumed` ist gesetzt. |
| `PGR-E5003` | Anfrage ohne aufgezeichnete Sitzung | Ihre Anwendung hat auf mehr Verbindungen Anfragen gestellt, als Sitzungen aufgezeichnet sind. Zeichnen Sie den Ablauf erneut auf, oder öffnen Sie weniger Verbindungen. |
| `PGR-E5004` | Antwort der Datenbank weicht von der Aufzeichnung ab | `play` mit `--compare-responses` hat bei einer Interaktion eine andere Antwort erhalten als aufgezeichnet. Die Meldung nennt Sitzung, Nummer der Interaktion, Position und Art der ersten Abweichung, bei der Art `Fehler` die SQLSTATE beider Seiten. Prüfen Sie, ob die Datenbank dieselbe Struktur hat wie die, gegen die Sie aufgezeichnet haben, oder zeichnen Sie erneut auf; siehe [Antworten beim Einspielen vergleichen](#antworten-beim-einspielen-vergleichen). Der Exit-Code ist 5. |
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
Ja. Eine Aufzeichnung im Format `yaml` ist Text und ändert sich nur, wenn sich die
Aufzeichnung ändert; zeichnen Sie mit `--record-timing` auf, können sich bei jedem
Aufzeichnen desselben Ablaufs die Werte von `offset_ms` ändern. Eine SQLite-Datei ist
binär. Beachten Sie den Hinweis zu vertraulichen Daten in der Einleitung.

**Kann ich die Datei von Hand bearbeiten?**
Das Werkzeug unterstützt das bei keinem der beiden Formate. Zeichnen Sie stattdessen
erneut auf. Eine SQLite-Aufzeichnung können Sie mit `sqlite3` abfragen (siehe
[Eine SQLite-Aufzeichnung abfragen](#eine-sqlite-aufzeichnung-abfragen)).

**Funktioniert die Wiedergabe, wenn ich die Anfragen umformuliere?**
Nein. Die Anfrage muss Zeichen für Zeichen der aufgezeichneten entsprechen,
auch in Leerzeichen, Kommentaren und Groß- und Kleinschreibung. Beim erweiterten
Protokoll gilt das auch für Namen, Parametertypen, Parameterwerte, Formate und
das Zeilenlimit. Ausgenommen sind nur Anfragen, die nur aus Leerraum und
Kommentaren bestehen (Lebendprüfungen, siehe
[Mit einem Datenbanktreiber arbeiten](#mit-einem-datenbanktreiber-arbeiten)).

**Kann ich Aufzeichnungen auf einem anderen Rechner verwenden?**
Ja, in beiden Formaten. Die Datei enthält keine rechnerspezifischen Angaben des
Werkzeugs.

## 9. Glossar

| Begriff | Bedeutung |
|---|---|
| Aufzeichnung | Datei mit der aufgezeichneten Kommunikation zwischen Anwendung und Datenbank, im Format `yaml` oder `sqlite` |
| Format | Art, wie die Datei die Aufzeichnung speichert: `yaml` (Text) oder `sqlite` (SQLite-Datei) |
| Verbindung | eine Datenbankverbindung Ihrer Anwendung; in der Aufzeichnung heißt sie Sitzung |
| Einfaches Anfrageprotokoll | PostgreSQL-Protokollvariante, in der jede Anfrage als einzelner Text gesendet wird |
| Erweitertes Anfrageprotokoll | Protokollvariante mit vorbereiteten Anweisungen und Parametern |
| Upstream | die echte Datenbank, an die `record` weiterleitet |
| Meldungscode | Kennung einer Meldung des Werkzeugs, zum Beispiel `PGR-E5001` |
| Interaktion | eine einfache Anfrage oder eine Folge des erweiterten Protokolls bis zu ihrem `Sync`, mit den Antworten der Datenbank; ihre Nummer in der Aufzeichnung steht in den Meldungen von `play` |
| Zeitangabe | das Feld `offset_ms` einer Anfrage in der Aufzeichnung: die Millisekunden seit der Startup-Nachricht der Verbindung bis zum Eintreffen der Anfrage, aufgezeichnet mit `--record-timing` |

## 10. Support und Kontakt

Melden Sie Fehler und Fragen über das Projekt-Repository. Geben Sie die
Programmversion, den Aufruf (ohne Zugangsdaten) und die Ausgabe mit
`--log-level debug` an. Entfernen Sie vertrauliche Daten, bevor Sie Ausgaben
oder Aufzeichnungen weitergeben.

## 11. Änderungshistorie

Es gibt keine veröffentlichte Version.
