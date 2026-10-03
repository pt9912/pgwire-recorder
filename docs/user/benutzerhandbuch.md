# Benutzerhandbuch: pgwire-recorder

Version: 0.1  
Software-Version: noch nicht veröffentlicht  
Stand: 03.10.2026

## 1. Einleitung

### Zweck der Software

`pgwire-recorder` zeichnet die Kommunikation zwischen Ihrer Anwendung und einer
PostgreSQL-Datenbank auf und spielt sie später ohne die Datenbank wieder ab.
Damit laufen Tests Ihrer Anwendung reproduzierbar, ohne dass dafür eine
PostgreSQL-Instanz bereitstehen muss.

Das Werkzeug arbeitet in zwei Betriebsarten:

* **Aufzeichnen (`record`):** Ihre Anwendung verbindet sich mit
  `pgwire-recorder` statt direkt mit PostgreSQL. Das Werkzeug leitet alles an
  die echte Datenbank weiter und speichert die Kommunikation in einer Datei.
* **Wiedergeben (`replay`):** Ihre Anwendung verbindet sich wieder mit
  `pgwire-recorder`. Das Werkzeug beantwortet die Anfragen aus der Datei. Eine
  Datenbank ist nicht nötig.

### Zielgruppe

Dieses Handbuch richtet sich an Entwicklerinnen, Entwickler und Tester, die
Anwendungen mit PostgreSQL-Anbindung testen, sowie an Teams, die diese Tests in
CI-Systemen ausführen.

### Voraussetzungen

* Ihre Anwendung kann Host und Port der Datenbankverbindung einstellen. Weitere
  Änderungen an der Anwendung sind nicht nötig.
* Zum Aufzeichnen: eine erreichbare PostgreSQL-Datenbank.
* Linux auf `amd64` oder `arm64`, oder ein Container-Laufzeitsystem.
* Die Verbindung zum Werkzeug läuft ohne Verschlüsselung (siehe
  [Fehlerbehebung](#7-fehlerbehebung)).

### Wichtig: Aufzeichnungen können vertrauliche Daten enthalten

Eine Aufzeichnung enthält die SQL-Anfragen und die zurückgegebenen Datensätze.
Das Werkzeug erkennt oder maskiert vertrauliche Inhalte nicht. Verwenden Sie
zum Aufzeichnen geeignete Testdaten, und speichern, übertragen und versionieren
Sie Aufzeichnungen so, wie es zu ihrem Inhalt passt.

## 2. Installation

### Binary

1. Laden Sie das Binary `pgwire-recorder` für Ihre Plattform herunter
   (`linux/amd64` oder `linux/arm64`).
2. Machen Sie die Datei ausführbar.
3. Prüfen Sie die Installation mit `pgwire-recorder version`.

### Container

Das Container-Image enthält das Binary und braucht keine weitere Software. Ein
Beispiel für eine Wiedergabe in Docker Compose:

```yaml
services:
  recorder:
    image: pgwire-recorder:latest
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
dem Ende jeder Verbindung und beim Beenden neu geschrieben.

#### Hinweise

* Existiert die Zieldatei bereits, bricht das Werkzeug ab (`PGR-E2002`). Wollen
  Sie sie ersetzen, ergänzen Sie `--force`.
* Verwenden Sie für einen Ablauf eine einzige Datenbankverbindung. Öffnet Ihre
  Anwendung mehrere Verbindungen, entsteht eine Aufzeichnung mit mehreren
  Verbindungen, die sich nicht wiedergeben lässt (siehe
  [Fehlerbehebung](#7-fehlerbehebung)).

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
  Anfrage an dieser Stelle entsprechen, Zeichen für Zeichen. Eine abweichende
  Anfrage wird als Fehler gemeldet; das Werkzeug liefert dann keine geratene
  Antwort (`PGR-E5001`).
* Jede neue Verbindung beginnt wieder am Anfang der Aufzeichnung.

### Mit einem Datenbanktreiber arbeiten

Das Werkzeug unterstützt das einfache und das erweiterte Anfrageprotokoll von
PostgreSQL. Viele Treiber nutzen standardmäßig das erweiterte Protokoll mit
vorbereiteten Anweisungen. Sie müssen am Treiber nur Host und Port ändern.

#### Hinweise

* Erzeugt ein Treiber bei jedem Lauf andere Namen für vorbereitete Anweisungen,
  passt die Wiedergabe nicht zur Aufzeichnung. Stellen Sie den Treiber dann so
  ein, dass er feste Namen verwendet oder keine Anweisungen zwischenspeichert.
* Die Verschlüsselung der Verbindung wird abgelehnt. Stellen Sie in Ihrer
  Verbindungszeichenfolge die Verschlüsselung ab, zum Beispiel mit
  `sslmode=disable`.

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

Alle Einstellungen lassen sich über Optionen und über Umgebungsvariablen
festlegen. Gilt dieselbe Einstellung mehrfach, setzt sich das Argument vor der
Umgebungsvariablen vor dem Standardwert durch.

| Option | Betriebsart | Umgebungsvariable | Standard |
|---|---|---|---|
| `--listen` | `record`, `replay` | `PGWIRE_RECORDER_LISTEN` | Pflicht |
| `--upstream` | `record` | `PGWIRE_RECORDER_UPSTREAM` | Pflicht |
| `--output` | `record` | `PGWIRE_RECORDER_OUTPUT` | Pflicht |
| `--force` | `record` | `PGWIRE_RECORDER_FORCE` | `false` |
| `--input` | `replay` | `PGWIRE_RECORDER_INPUT` | Pflicht |
| `--fail-on-unconsumed` | `replay` | `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` | `false` |
| `--log-level` | `record`, `replay` | `PGWIRE_RECORDER_LOG_LEVEL` | `info` |

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
| 5 | Abweichung bei der Wiedergabe |
| 6 | nicht unterstützte Funktion des Protokolls |

Ein Fehler, der nur eine Verbindung betrifft, beendet diese Verbindung. Das
Werkzeug läuft weiter und liefert den Exit-Code des ersten aufgetretenen Fehlers
erst, wenn Sie es beenden.

## 6. Rollen und Rechte

Das Werkzeug kennt keine Benutzer, Rollen oder Anmeldung. Es prüft keine
Zugangsdaten selbst: Beim Aufzeichnen leitet es die Anmeldung an die Datenbank
weiter, beim Wiedergeben nimmt es jede Anmeldung an. Betreiben Sie es nur in
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
| `PGR-E2000`, `PGR-E2001` | ungültiger Aufruf oder ungültige Konfiguration | Eine Option fehlt, ist unbekannt oder hat einen ungültigen Wert. Prüfen Sie den Aufruf mit `--help`. |
| `PGR-E2002` | Zieldatei existiert bereits | Wählen Sie einen anderen Dateinamen, oder ergänzen Sie `--force`, um die Datei zu ersetzen. |
| `PGR-E3000`, `PGR-E3001` | Aufzeichnung nicht lesbar oder nicht schreibbar | Die Datei fehlt, oder Sie haben keine Rechte. Prüfen Sie Pfad und Dateirechte. |
| `PGR-E3002` | unbekannte Version der Aufzeichnung | Die Datei stammt aus einer anderen Programmversion. Zeichnen Sie mit der verwendeten Version erneut auf. |
| `PGR-E3003` | Aufzeichnung beschädigt | Die Datei ist unvollständig oder verändert. Zeichnen Sie erneut auf. |
| `PGR-E3004` | Aufzeichnung ohne Verbindung | Beim Aufzeichnen hat sich keine Anwendung verbunden. Zeichnen Sie erneut auf. |
| `PGR-E4000`, `PGR-E4003` | Verbindung unerwartet beendet | Die Verbindung brach mitten in einer Anfrage ab. Prüfen Sie Netzwerk, Datenbank und Anwendung. |
| `PGR-E4001` | Adresse nicht nutzbar | Der Port aus `--listen` ist belegt oder nicht erlaubt. Wählen Sie einen freien Port. |
| `PGR-E4002` | Datenbank nicht erreichbar | Prüfen Sie `--upstream`, die Datenbank und das Netzwerk. |
| `PGR-E5000`, `PGR-E5001` | Abweichung bei der Wiedergabe | Ihre Anwendung hat eine andere Anfrage gestellt als aufgezeichnet. Die Meldung nennt die erwartete und die empfangene Anfrage. Zeichnen Sie erneut auf, oder korrigieren Sie die Anwendung. |
| `PGR-E5002` | aufgezeichnete Anfragen nicht verbraucht | Ihr Test hat weniger Anfragen gestellt als aufgezeichnet, und `--fail-on-unconsumed` ist gesetzt. |
| `PGR-E6000`, `PGR-E6001` | nicht unterstützte Nachricht | Die Anwendung nutzt eine Funktion, die das Werkzeug nicht unterstützt, zum Beispiel `COPY`. Verwenden Sie diese Funktion im aufgezeichneten Ablauf nicht. |
| `PGR-E6002` | nicht unterstützte Protokollversion | Das Werkzeug unterstützt Version 3.0 des PostgreSQL-Protokolls. Verwenden Sie einen Treiber, der sie nutzt. |
| `PGR-E6003` | Aufzeichnung mit mehreren Verbindungen | Eine Aufzeichnung mit mehr als einer Verbindung lässt sich nicht wiedergeben. Zeichnen Sie den Ablauf mit einer einzigen Verbindung auf, zum Beispiel mit einer Pool-Größe von 1. |

### Warnungen

| Code | Bedeutung | Hinweis |
|---|---|---|
| `PGR-W2001` | Wiedergabe endete vor der letzten aufgezeichneten Anfrage | Ihr Test hat nicht alle aufgezeichneten Anfragen ausgeführt. Mit `--fail-on-unconsumed` wird das zum Fehler. |
| `PGR-W3001` | Abbruchwunsch nicht weitergeleitet | Die Anwendung hat versucht, eine laufende Anfrage abzubrechen. Das Werkzeug leitet diesen Wunsch nicht weiter und schließt die Verbindung. |

### Die Anwendung kann sich nicht verbinden

Fordert Ihr Treiber eine verschlüsselte Verbindung an und bricht ab, wenn sie
abgelehnt wird, schalten Sie die Verschlüsselung in der Verbindungszeichenfolge
ab (zum Beispiel `sslmode=disable`).

## 8. FAQ

**Brauche ich für die Wiedergabe eine Datenbank?**
Nein. Die Wiedergabe beantwortet alle aufgezeichneten Anfragen aus der Datei.

**Kann ich eine Aufzeichnung in die Versionsverwaltung legen?**
Ja. Die Datei ist Text (YAML) und ändert sich nur, wenn sich die Aufzeichnung
ändert. Beachten Sie den Hinweis zu vertraulichen Daten in der Einleitung.

**Kann ich die Datei von Hand bearbeiten?**
Das Werkzeug unterstützt das nicht. Zeichnen Sie stattdessen erneut auf.

**Funktioniert die Wiedergabe, wenn ich die Anfragen umformuliere?**
Nein. Die Anfrage muss Zeichen für Zeichen der aufgezeichneten entsprechen,
auch in Leerzeichen, Kommentaren und Groß- und Kleinschreibung.

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
