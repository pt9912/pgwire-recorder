# Lastenheft — pgwire-recorder

**Version:** 0.1.0 (SemVer 2.0, `Major.Minor.Patch`). Ab `Accepted` ist der Bump der
Fußabdruck des Change Requests: **Major** bei Änderung oder Streichung einer
bestehenden Anforderung, **Minor** bei einer neuen Anforderung, **Patch** bei
einer Tatsachenberichtigung ohne Änderung einer Aussage (Baseline-Regelwerk
`grundlagen-source-precedence.md` §Spec-Stratifizierung).

**Status:** Draft — gilt dem **Dokument**, nicht der
einzelnen Anforderung. Vor `Accepted` frei änderbar, ohne Change Request, ohne
Historie-Zeile; ab `Accepted` ist jede Änderung eine Vertragsänderung
(Baseline-Regelwerk `grundlagen-source-precedence.md`
§Spec-Stratifizierung).

**Autor:** pt9912, **Datum:** 2026-10-03.

**Zielversion des Produkts:** v1. Das Produkt ist fertig, wenn alle Anforderungen dieses Lastenhefts (MUSS und SOLL) umgesetzt sind.

---

## 1. Zweck und Geltungsbereich

Dieses Lastenheft beschreibt die fachlichen Anforderungen an `pgwire-recorder`
aus Sicht der Anwender und Auftraggeber. Es beschreibt, **was** das Produkt
leisten soll und welche Ziele damit verfolgt werden. Technische Entscheidungen
zu Implementierung, Programmiersprache, Bibliotheken, internen Datenstrukturen
und Softwarearchitektur gehören nicht in dieses Dokument; sie werden in den
nachgelagerten Dokumenten festgelegt.

**Ausgangssituation.** Automatisierte Tests von Anwendungen mit
PostgreSQL-Abhängigkeiten benötigen häufig eine reale oder eigens bereitgestellte
PostgreSQL-Instanz. Das erhöht den Aufwand für lokale Entwicklungsumgebungen und
CI-Systeme und kann Tests langsamer, komplexer und weniger deterministisch
machen. Für viele Tests ist nicht die Datenbank selbst Gegenstand des Tests,
sondern das aus Sicht der Anwendung beobachtbare Verhalten der
PostgreSQL-Kommunikation.

**Produktziel.** `pgwire-recorder` ist ein eigenständig ausführbares
Kommandozeilenwerkzeug, das als Vermittler zwischen PostgreSQL-Clients und
PostgreSQL-Servern eingesetzt wird. Es bietet drei Betriebsarten:

- **Record:** Die Anwendung verbindet sich mit `pgwire-recorder` statt direkt
  mit PostgreSQL. Der Recorder leitet die unterstützte Kommunikation an einen
  realen PostgreSQL-Server weiter und zeichnet die für eine spätere Wiedergabe
  erforderlichen Interaktionen auf.
- **Replay:** Die Anwendung verbindet sich ebenfalls mit `pgwire-recorder`. Für
  die aufgezeichneten und unterstützten Interaktionen ist kein realer
  PostgreSQL-Server erforderlich; der Recorder beantwortet die Anfragen anhand
  einer zuvor erzeugten Aufzeichnung.
- **Einspielen:** Der Recorder führt die in einer Aufzeichnung enthaltenen
  Client-Anfragen gegen einen PostgreSQL-Server aus, ohne dass eine Anwendung
  beteiligt ist.

```text
Record:  Application --PG Wire--> pgwire-recorder --> PostgreSQL
                                        |
                                        +--> Recording

Replay:  Application --PG Wire--> pgwire-recorder <-- Recording
                                  (kein PostgreSQL-Server erforderlich)

Einspielen:  Recording --> pgwire-recorder --PG Wire--> PostgreSQL
                           (keine Anwendung beteiligt)
```

**Produktgrenzen.** `pgwire-recorder` ist kein vollständiger Ersatz für
PostgreSQL. Das Produkt simuliert nicht die interne Datenbanklogik,
SQL-Ausführung, Transaktionen oder Datenhaltung. Im Replay-Modus wird
ausschließlich Verhalten reproduziert, das durch eine geeignete Aufzeichnung
abgedeckt und vom jeweiligen Produktstand unterstützt wird. Das Produkt ist
ebenfalls kein allgemeiner Netzwerk-Paketmitschnitt und kein
Datenbankadministrationswerkzeug; das Einspielen führt ausschließlich
aufgezeichnete Anfragen aus.

**Primäre Anwendungsfälle.**

- *Datenbankkommunikation aufzeichnen:* Ein Entwickler startet einen realen
  PostgreSQL-Server und `pgwire-recorder` im Record-Modus, die zu testende
  Anwendung verbindet sich mit dem Recorder; danach steht eine persistente
  Aufzeichnung für spätere Testläufe zur Verfügung.
- *Aufzeichnung wiedergeben:* Ein Entwickler oder ein CI-System startet
  `pgwire-recorder` im Replay-Modus mit einer vorhandenen Aufzeichnung; die
  Anwendung erhält die zuvor aufgezeichneten Antworten.
- *Aufzeichnung in eine Datenbank einspielen:* Ein Entwickler oder ein
  CI-System startet `pgwire-recorder` im Einspiel-Modus mit einer vorhandenen
  Aufzeichnung und einem PostgreSQL-Server; die aufgezeichneten Anfragen werden
  gegen diesen Server ausgeführt. Damit lässt sich eine aufgezeichnete Last
  reproduzierbar erzeugen, zum Beispiel um eine Komponente zu testen, die
  Änderungen der Datenbank verarbeitet (Change Data Capture).
- *Aufzeichnung versionieren:* Eine Aufzeichnung liegt als Datei oder als
  zusammengehöriger Satz von Dateien vor und kann gemeinsam mit Testcode
  gespeichert, transportiert und — sofern vom Anwender gewünscht — versioniert
  werden.

## 2. Stakeholder

| Stakeholder | Rolle | Erwartung |
|---|---|---|
| Softwareentwickler | Anwender | PostgreSQL-abhängige Anwendungen ohne Datenbankinstanz testen können |
| Entwickler und Tester | Anwender | Reproduzierbare Integrations- und Komponententests |
| Teams mit CI-Umgebungen | Anwender | PostgreSQL-abhängige Tests nicht-interaktiv und ohne Datenbankdienst ausführen |
| Maintainer von Testinfrastrukturen | Betreiber | Einfach bereitstellbares, automatisierbares Werkzeug |

## 3. Funktionale Anforderungen

Regeln dieser Sektion: ID-Schema `LH-FA-<NN>`. Jede Anforderung trägt drei
Pfade — Happy · Boundary · Negative — plus Out-of-Scope (Baseline-Regelwerk
`modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium). Eine Anforderung, deren
Bedarf ersatzlos entfällt, wird **nicht gelöscht**, sondern trägt den Vermerk
*zurückgezogen* im Titel; die Nummer bleibt vergeben (Baseline-Regelwerk
`grundlagen-source-precedence.md` §Spec-Stratifizierung).

**Priorität:** **MUSS** = erforderlich für v1 · **SOLL** = erwünscht für v1,
sofern ohne unverhältnismäßigen Aufwand realisierbar · **KANN** = mögliche
spätere Erweiterung. Fertig ist das Produkt erst, wenn auch die SOLL-Anforderungen
umgesetzt sind.

**Globaler Funktionsumfang v1:** PostgreSQL-Kommunikation über das Simple Query
Protocol und das Extended Query Protocol. Anfragen und Interaktionen außerhalb
dieses Umfangs gelten als „nicht unterstützt" (siehe LH-FA-05, LH-FA-18 und §5).

### LH-FA-01 — Kommandozeilenanwendung

**Priorität:** MUSS

**Beschreibung:** `pgwire-recorder` muss als eigenständig ausführbares
Kommandozeilenwerkzeug nutzbar sein.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine bereitgestellte Installation, when das Werkzeug von
  der Kommandozeile gestartet wird, then läuft es ohne weitere Software des
  Produkts (kein GUI, kein zentraler Server).
- **Boundary:** Given ein Start ohne gültige Betriebsart, when das Werkzeug
  aufgerufen wird, then endet es mit einer erkennbaren Fehlermeldung statt in
  einem undefinierten Zustand.
- **Negative:** Given ungültige Aufrufparameter, when das Werkzeug gestartet
  wird, then wird ein Konfigurationsfehler signalisiert (siehe LH-FA-13).

**Out-of-Scope:** Grafische Benutzeroberfläche; zentraler Recording-Server.

---

### LH-FA-02 — Record-Modus

**Priorität:** MUSS

**Beschreibung:** Das Produkt muss einen Record-Modus bereitstellen, in dem
unterstützte PostgreSQL-Kommunikation zwischen Client und realem
PostgreSQL-Server vermittelt und aufgezeichnet wird.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein erreichbarer PostgreSQL-Server und ein Client, der
  sich mit dem Recorder verbindet, when der Client eine unterstützte
  Simple-Query-Interaktion ausführt, then erhält der Client die Serverantwort
  und eine Aufzeichnung entsteht (siehe Abnahmeszenario 1 in §7).
- **Boundary:** Given ein Client, der keine Anfrage stellt, when die Verbindung
  endet, then wird diese Verbindung nicht aufgezeichnet, die
  Aufzeichnung bleibt gültig; der Recorder bleibt nicht hängen.
- **Negative:** Given ein nicht erreichbarer PostgreSQL-Server, when ein Client
  sich verbindet, then wird ein Verbindungsfehler signalisiert (siehe LH-QA-05).

**Out-of-Scope:** Veränderung von Serverantworten.

---

### LH-FA-03 — Replay-Modus

**Priorität:** MUSS

**Beschreibung:** Das Produkt muss einen Replay-Modus bereitstellen, in dem eine
zuvor erzeugte Aufzeichnung ohne realen PostgreSQL-Server wiedergegeben werden
kann.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine vorhandene Aufzeichnung und kein erreichbarer
  PostgreSQL-Server, when ein Client dieselbe unterstützte Interaktion ausführt,
  then erhält er das aufgezeichnete, aus Sicht des Clients relevante Verhalten
  (siehe Abnahmeszenario 2 in §7).
- **Boundary:** Given ein Replay ohne jede Netzverbindung zu PostgreSQL, when
  das Werkzeug startet, then versucht es keine Verbindung zu einem Server.
- **Negative:** Given eine ungültige oder nicht lesbare Aufzeichnung, when das
  Replay gestartet wird, then wird dies als ungültiges Recording signalisiert
  (siehe LH-QA-05).

**Out-of-Scope:** Simulation frei definierter PostgreSQL-Antworten ohne
vorheriges Recording.

---

### LH-FA-04 — Transparenz für Clients

**Priorität:** MUSS

**Beschreibung:** Bestehende PostgreSQL-Clients sollen für die Nutzung
grundsätzlich keine programmspezifische Integration des Recorders benötigen. Für
einen unterstützten Anwendungsfall soll es genügen, Verbindungsadresse
beziehungsweise Host und Port auf `pgwire-recorder` umzustellen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein unterstützter Testfall, when Host und Port auf den
  Recorder umgestellt werden, then läuft der Testfall ohne Änderung der
  fachlichen Programmlogik (siehe Abnahmeszenario 3 in §7).
- **Boundary:** Given ein Client, der Funktionen außerhalb des unterstützten
  Umfangs nutzt, when er sich verbindet, then wird die Nichtunterstützung
  erkennbar gemeldet, statt stillschweigend falsches Verhalten zu zeigen.
- **Negative:** Given ein Client, der zwingend Recorder-spezifische
  Programmlogik benötigte, then gilt der Anwendungsfall als nicht unterstützt.

**Out-of-Scope:** Anpassung von Client-Bibliotheken.

---

### LH-FA-05 — Simple Query Protocol

**Priorität:** MUSS

**Beschreibung:** v1 muss PostgreSQL-Kommunikation über das Simple Query
Protocol unterstützen. Der Funktionsumfang von v1 umfasst den Verbindungsaufbau
(Start-up einschließlich Authentifizierung), Anfragen (`Query`), das
Verbindungsende (`Terminate`) und die Serverantworten einer Anfrage bis
`ReadyForQuery`, einschließlich Fehler- und Hinweisantworten sowie
Transaktionsbefehlen als gewöhnliche Anfragen. TLS-Aushandlung ist nur als
Verschlüsselung zum Client nach LH-FA-23 Teil des Funktionsumfangs.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Simple-Query-Anfrage im unterstützten Umfang, when
  sie über den Recorder läuft, then wird sie aufgezeichnet beziehungsweise
  wiedergegeben.
- **Boundary:** Given eine Interaktion am Rand des unterstützten Umfangs, when
  sie auftritt (etwa eine TLS-Anfrage ohne Verschlüsselung nach LH-FA-23, ein
  Abbruchwunsch oder eine
  Datenübertragung per `COPY`), then ist ihr Status (unterstützt / nicht
  unterstützt) eindeutig und dokumentiert; eine nicht unterstützte Interaktion
  wird nicht still ignoriert.
- **Negative:** Given eine nicht unterstützte Protokollinteraktion, when sie
  auftritt, then wird sie mit einem klar erkennbaren Fehler abgelehnt (siehe
  LH-QA-05).

**Out-of-Scope:** — (das Extended Query Protocol regelt LH-FA-18).

---

### LH-FA-06 — Aufzeichnung von Anfragen und Antworten

**Priorität:** MUSS

**Beschreibung:** Eine Aufzeichnung muss genügend Informationen über
unterstützte Client-Anfragen und zugehörige Serverantworten enthalten, um diese
im Replay-Modus reproduzieren zu können.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine aufgezeichnete unterstützte Interaktion, when sie
  im Replay angefragt wird, then genügt die Aufzeichnung allein zur
  Beantwortung.
- **Boundary:** Given eine Serverantwort mit mehreren Ergebnismengen oder ohne
  Datensätze, when sie aufgezeichnet wird, then bleibt sie vollständig
  reproduzierbar.
- **Negative:** Given eine nicht unterstützte Interaktion, when sie auftritt,
  then wird sie nicht als scheinbar vollständige Aufzeichnung abgelegt.

**Out-of-Scope:** Automatische Maskierung oder Redaktion sensibler Inhalte
(siehe LH-RB-01).

---

### LH-FA-07 — Persistente Recordings

**Priorität:** MUSS

**Beschreibung:** Aufzeichnungen müssen persistent gespeichert und bei einem
späteren Programmstart erneut verwendet werden können. Eine Aufzeichnung liegt
als Datei oder als zusammengehöriger Satz von Dateien vor und kann gemeinsam mit
Testcode gespeichert, transportiert und versioniert werden.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine abgeschlossene Aufzeichnung, when das Werkzeug neu
  gestartet wird, then ist sie für ein Replay verwendbar.
- **Boundary:** Given eine Aufzeichnung, die in ein anderes Verzeichnis oder auf
  einen anderen Rechner kopiert wurde, when sie dort verwendet wird, then
  funktioniert das Replay unverändert; der Recorder trägt keine
  rechnerspezifischen Angaben in die Aufzeichnung ein.
- **Negative:** Given eine unvollständige oder beschädigte Datei, when sie
  geladen wird, then wird sie als ungültiges Recording erkannt.

**Out-of-Scope:** Zentrales Recording-Repository.

---

### LH-FA-08 — Auswahl eines Recordings

**Priorität:** MUSS

**Beschreibung:** Anwender müssen bestimmen können, welche Aufzeichnung erzeugt
beziehungsweise für ein Replay verwendet wird.

**Akzeptanzkriterien:**

- **Happy Path:** Given mehrere Aufzeichnungen, when der Anwender eine davon
  benennt, then wird genau diese erzeugt beziehungsweise wiedergegeben.
- **Boundary:** Given ein Aufnahmeziel, das bereits existiert, when Record
  gestartet wird, then wird das Ziel nicht stillschweigend überschrieben; der
  Anwender kann das Überschreiben ausdrücklich verlangen.
- **Negative:** Given eine nicht existierende Aufzeichnung für Replay, when das
  Werkzeug startet, then wird ein Fehler signalisiert.

**Out-of-Scope:** Verwaltung eines zentralen Recording-Repositorys.

---

### LH-FA-09 — Reproduzierbares Replay

**Priorität:** MUSS

**Beschreibung:** Bei einer zum Recording passenden unterstützten Anfrage muss
der Replay-Modus das aufgezeichnete Verhalten reproduzierbar bereitstellen.

**Akzeptanzkriterien:**

- **Happy Path:** Given dieselbe Aufzeichnung und dieselbe unterstützte
  Anfrage, when sie mehrfach wiedergegeben wird, then ist das beobachtbare
  Verhalten jedes Mal gleich.
- **Boundary:** Given mehrfach identische Queries in einer Aufzeichnung, when
  sie wiedergegeben werden, then erhält jede Anfrage die zu ihrer Position in der
  Aufzeichnung gehörende Antwort.
- **Negative:** Given eine Anfrage ohne passende Aufzeichnung, when sie
  eintrifft, then greift LH-FA-10.

**Out-of-Scope:** Veränderung oder Manipulation aufgezeichneter
Datenbankantworten.

---

### LH-FA-10 — Abweichende Anfrage

**Priorität:** MUSS

**Beschreibung:** Kann eine Anfrage im Replay-Modus nicht anhand der
Aufzeichnung bedient werden, muss dies als klar erkennbarer Fehler behandelt
werden. Der Recorder darf in diesem Fall nicht stillschweigend eine beliebige
oder offensichtlich unpassende Antwort verwenden.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine unterstützte Anfrage, die nicht zur Aufzeichnung
  passt, when sie im Replay eintrifft, then liefert beziehungsweise signalisiert
  der Recorder einen eindeutigen Fehler (siehe Abnahmeszenario 4 in §7).
- **Boundary:** Given eine Anfrage, die einer aufgezeichneten nur ähnlich ist,
  when sie eintrifft, then wird sie nur bedient, wenn sie der aufgezeichneten
  Anfrage entspricht, andernfalls als Abweichung gemeldet — nie mit einer
  geratenen Antwort.
- **Negative:** Given eine Abweichung, when sie auftritt, then wird keine
  aufgezeichnete Antwort einer anderen Anfrage geliefert.

**Out-of-Scope:** Automatisches Nachrecorden fehlender Interaktionen.

---

### LH-FA-11 — Fehler des PostgreSQL-Servers

**Priorität:** MUSS

**Beschreibung:** Unterstützte PostgreSQL-Fehlerantworten, die während einer
Aufzeichnung auftreten, müssen so aufgezeichnet werden können, dass das
entsprechende beobachtbare Fehlerverhalten beim Replay reproduziert werden kann.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Anfrage, die im Record eine unterstützte
  Fehlerantwort erzeugt, when sie im Replay erneut gestellt wird, then erhält der
  Client ein äquivalentes beobachtbares Fehlerverhalten (siehe Abnahmeszenario 6 in §7).
- **Boundary:** Given eine Fehlerantwort innerhalb einer Folge mehrerer
  Interaktionen, when sie wiedergegeben wird, then bleibt ihre Position in der
  Reihenfolge erhalten (siehe LH-FA-12).
- **Negative:** Given eine Fehlerantwort, die nicht verlustfrei aufgezeichnet
  werden kann, when sie im Record auftritt, then wird dies als nicht
  unterstützte Protokollinteraktion gemeldet, statt eine unvollständige
  Aufzeichnung abzulegen.

**Out-of-Scope:** Erzeugen frei definierter Fehlerantworten ohne vorheriges
Recording.

---

### LH-FA-12 — Geordnete Interaktionen

**Priorität:** MUSS

**Beschreibung:** Das Recording muss die für die Wiedergabe relevante Reihenfolge
der aufgezeichneten Interaktionen erhalten.

**Akzeptanzkriterien:**

- **Happy Path:** Given mehrere nacheinander aufgezeichnete Interaktionen, when
  sie wiedergegeben werden, then entspricht die Reihenfolge der aufgezeichneten.
- **Boundary:** Given mehrere Client-Verbindungen, when sie aufgezeichnet
  werden, then wird jede Verbindung mit Anfragen als eigene Session
  aufgezeichnet, und die Reihenfolge der Interaktionen bleibt je Session
  erhalten; auf Wunsch werden auch Verbindungen ohne Anfrage als Session
  aufgezeichnet. Given ein Recording mit mehreren Sessions, when Clients
  nacheinander Verbindungen aufbauen, then erhält die n-te Verbindung mit einer
  Anfrage die n-te aufgezeichnete Session; auf Wunsch zählt stattdessen die n-te
  Verbindung, was eine Aufzeichnung mit Verbindungen ohne Anfrage voraussetzt.
- **Negative:** Given eine Folge, die der aufgezeichneten Reihenfolge
  widerspricht und für die Reihenfolge relevant ist, when sie im Replay
  eintrifft, then wird sie nach LH-FA-10 behandelt.

**Out-of-Scope:** Eine deterministische Zuordnung bei gleichzeitigem
Verbindungsaufbau.

---

### LH-FA-13 — Prozessbeendigung und Fehlerstatus

**Priorität:** MUSS

**Beschreibung:** Das CLI-Tool muss Betriebs- und Konfigurationsfehler für
automatisierte Umgebungen erkennbar signalisieren.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein fehlerfreier Lauf, when das Werkzeug endet, then
  signalisiert der Prozessstatus Erfolg.
- **Boundary:** Given ein kontrolliertes Beenden durch Signal, when es
  eintritt, then signalisiert der Prozessstatus Erfolg, wenn bis dahin kein
  Fehler aufgetreten ist, andernfalls den Fehler.
- **Negative:** Given ein Betriebs- oder Konfigurationsfehler, when er auftritt,
  then wird er mit seiner Fehlerklasse gemeldet, und der
  Prozessstatus signalisiert Misserfolg, sobald der Prozess endet; er
  unterscheidet die Fehlerklassen aus LH-QA-05.

**Out-of-Scope:** Eine bestimmte Wertebelegung der Exit Codes.

---

### LH-FA-14 — Diagnoseausgaben

**Priorität:** SOLL

**Beschreibung:** Das Tool soll verständliche Diagnoseinformationen ausgeben,
insbesondere bei Verbindungsproblemen, ungültigen Aufzeichnungen und
Replay-Abweichungen.

**Akzeptanzkriterien:**

- **Happy Path:** Given eines der genannten Problemszenarien, when es eintritt,
  then nennt die Ausgabe Art und Ursache des Problems.
- **Boundary:** Given ein Lauf ohne Probleme, when er endet, then bleibt die
  Ausgabe auf der konfigurierten Detailstufe; die Detailstufe ist einstellbar.
- **Negative:** Given ein Fehler, when er gemeldet wird, then enthält die Meldung
  keine irreführende Fehlerklasse.

**Out-of-Scope:** Grafische Auswertung.

---

### LH-FA-15 — CI-Eignung

**Priorität:** MUSS

**Beschreibung:** Das Tool muss ohne interaktive Benutzereingaben ausführbar
sein, damit es in automatisierten Tests und CI-Systemen eingesetzt werden kann.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein CI-Lauf ohne Terminal, when Record oder Replay
  gestartet wird, then verläuft der Lauf ohne jede Eingabeaufforderung (siehe
  Abnahmeszenario 5 in §7).
- **Boundary:** Given fehlende Konfiguration, when das Werkzeug startet, then
  bricht es mit Fehler ab, statt nachzufragen.
- **Negative:** Given eine Situation, die sonst eine Rückfrage auslösen würde,
  when sie eintritt, then ist das Ergebnis ein Fehler, kein Warten auf Eingabe.

**Out-of-Scope:** Interaktive Betriebsarten.

---

### LH-FA-16 — Container-Eignung

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll als Docker/OCI-Image bereitgestellt werden
und so betreibbar sein, dass es ohne besondere Anforderungen in
containerisierten Testumgebungen eingesetzt werden kann.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine containerisierte Testumgebung, when das
  bereitgestellte Image dort gestartet wird, then funktionieren Record und
  Replay ohne zusätzliche Sonderkonfiguration.
- **Boundary:** Given eine Umgebung ohne Terminal, when das Werkzeug startet,
  then ist das unproblematisch (siehe LH-FA-15).
- **Negative:** Given eine Umgebung ohne Netzzugang zu PostgreSQL, when Replay
  läuft, then funktioniert es (siehe LH-FA-03).

**Out-of-Scope:** Bereitstellung orchestrierter Deployments; Windows-Container.

---

### LH-FA-17 — Maschinenlesbare Konfiguration

**Priorität:** SOLL

**Beschreibung:** Alle für automatisierte Testläufe erforderlichen Einstellungen
sollen über nicht-interaktive Mechanismen festgelegt werden können: über
Argumente, Umgebungsvariablen und eine Konfigurationsdatei, in der sich auch
Verbindungen zu Datenbanken benennen lassen. Die gewählte Konfiguration lässt sich
anzeigen, ohne Geheimnisse preiszugeben.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine nicht-interaktive Konfiguration, when das Werkzeug
  startet, then sind alle für Record, Replay und Einspielen nötigen
  Einstellungen darüber setzbar.
- **Boundary:** Given dieselbe Einstellung aus mehreren Quellen, when das
  Werkzeug startet, then gilt eine eindeutige, dokumentierte Priorität; Geheimnisse
  wie ein Passwort stehen nicht im Klartext in der Konfigurationsdatei, sondern
  werden aus der Umgebung eingesetzt, und die Anzeige der Konfiguration zeigt sie nicht.
- **Negative:** Given eine ungültige Konfiguration, when das Werkzeug startet,
  then wird ein Konfigurationsfehler signalisiert (siehe LH-FA-13).

**Out-of-Scope:** Interaktive Konfigurationsassistenten.

---

### LH-FA-18 — Extended Query Protocol

**Priorität:** MUSS

**Beschreibung:** Das Produkt muss PostgreSQL-Kommunikation über das Extended
Query Protocol (`Parse`, `Bind`, `Describe`, `Execute`, `Sync`, Prepared
Statements) im Record- und im Replay-Modus unterstützen. Damit genügt für
Clients, die dieses Protokoll standardmäßig verwenden, die Umstellung von Host
und Port (LH-FA-04).

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Client, der das Extended Query Protocol verwendet,
  when seine Interaktion über den Recorder aufgezeichnet und später im Replay
  wiedergegeben wird, then erhält er das aufgezeichnete, aus Sicht des Clients
  relevante Verhalten (siehe Abnahmeszenario 7 in §7).
- **Boundary:** Given dasselbe Prepared Statement, das mehrfach mit
  unterschiedlichen Parameterwerten ausgeführt wird, when es im Replay
  wiedergegeben wird, then wird jede Ausführung ihrer aufgezeichneten Position
  zugeordnet.
- **Negative:** Given eine Extended-Query-Interaktion, die nicht zur
  Aufzeichnung passt, when sie im Replay eintrifft, then wird sie wie in
  LH-FA-10 als Abweichung gemeldet.

**Out-of-Scope:** `COPY`; Replikationsprotokoll.

---

### LH-FA-19 — Bereitstellung über Homebrew

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll sich unter macOS und Linux über den
Paketmanager Homebrew installieren lassen.

**Akzeptanzkriterien:**

- **Happy Path:** Given ein Rechner mit Homebrew unter macOS oder Linux, when der
  Anwender die dokumentierten Installationsbefehle ausführt, then ist
  `pgwire-recorder` ausführbar und meldet seine Version (siehe
  Abnahmeszenario 11 in §7).
- **Boundary:** Given eine Vorabversion, when sie veröffentlicht wird, then
  ändert sie die über Homebrew bereitgestellte Version nicht.
- **Negative:** Given Windows, auf dem Homebrew nicht verfügbar ist, when der
  Anwender das Produkt installieren will, then steht das Binary direkt zur
  Verfügung (siehe LH-FA-01).

**Out-of-Scope:** Aufnahme in das Standard-Repository von Homebrew.

---

### LH-FA-20 — Einspielen einer Aufzeichnung

**Priorität:** MUSS

**Beschreibung:** Das Produkt muss die in einer Aufzeichnung enthaltenen
Client-Anfragen gegen einen PostgreSQL-Server ausführen können (Einspielen),
ohne dass eine Anwendung beteiligt ist. Das Einspielen führt ausschließlich
aufgezeichnete Anfragen aus. Es authentifiziert sich gegenüber dem Server und
verbindet sich auf Wunsch verschlüsselt; dabei prüft es das Zertifikat des Servers
und lässt sich eine eigene Zertifizierungsstelle des Anwenders hinterlegen. Ein
Überspringen der Prüfung gibt es nicht.

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Aufzeichnung und ein erreichbarer PostgreSQL-Server,
  when das Einspielen gestartet wird, then werden die aufgezeichneten Anfragen in
  der aufgezeichneten Reihenfolge ausgeführt, und der Server enthält danach die
  Wirkung der Anfragen (siehe Abnahmeszenario 12 in §7).
- **Boundary:** Given eine Aufzeichnung mit mehreren Sessions, when sie
  eingespielt wird, then wird jede Session über eine eigene Verbindung
  ausgeführt, die Sessions nacheinander in der aufgezeichneten Reihenfolge.
  Given ein Fehler des Servers bei einer Anfrage, die auch in der Aufzeichnung
  mit einem Fehler beantwortet wurde, when der Anwender das ausdrücklich wählt,
  then gilt dieser Fehler als erwartet. Given ein Abbruchsignal, when es
  eintrifft, then endet das Einspielen kontrolliert; auf Wunsch wird zuvor die
  laufende Session beendet.
- **Negative:** Given ein nicht erreichbarer Server, eine fehlgeschlagene
  Authentifizierung oder eine Anfrage, die der Server mit einem unerwarteten
  Fehler beantwortet, when das Einspielen läuft, then wird der Fehler eindeutig
  gemeldet, und das Einspielen bricht ab, sofern der Anwender nicht ausdrücklich
  verlangt weiterzulaufen (siehe LH-QA-05).

**Out-of-Scope:** Vergleich der Serverantworten mit der Aufzeichnung (nur auf
Wunsch, siehe LH-FA-24); paralleles Einspielen mehrerer Sessions; zeitgetreues
Abspielen (siehe LH-FA-21, nur auf Wunsch).

---

### LH-FA-21 — Zeitgetreues Einspielen

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll beim Aufzeichnen auf Wunsch den zeitlichen
Abstand der Anfragen festhalten und beim Einspielen auf Wunsch wiederherstellen,
wahlweise relativ zur jeweils vorigen Anfrage oder absolut ab einem Bezugspunkt
(Verbindungsaufbau oder erste Anfrage). Ohne diese Wünsche enthält eine
Aufzeichnung keine Zeitangaben, und das Einspielen führt die Anfragen
nacheinander und ohne Wartezeiten aus (LH-FA-20).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Aufzeichnung mit Zeitangaben, when der Anwender
  zeitgetreues Einspielen verlangt, then ist bei der relativen Wiedergabe der
  Abstand zwischen zwei aufeinander folgenden Anfragen nicht kürzer als
  aufgezeichnet, und bei der absoluten Wiedergabe liegt keine Anfrage früher als im
  aufgezeichneten Abstand zum Bezugspunkt (siehe Abnahmeszenario 13 in §7).
- **Boundary:** Given keine Zeitaufzeichnung verlangt, when aufgezeichnet wird,
  then enthält die Aufzeichnung keine Zeitangaben; given keine zeitgetreue
  Wiedergabe verlangt, when eine Aufzeichnung mit Zeitangaben eingespielt wird,
  then werden die Zeitangaben ignoriert; given die absolute Wiedergabe, when der
  Bezugspunkt gewählt wird, then gilt der Verbindungsaufbau oder die erste
  Anfrage.
- **Negative:** Given eine Aufzeichnung ohne Zeitangaben, when zeitgetreues
  Einspielen verlangt wird, then wird dies als ungültige Verwendung gemeldet.

**Out-of-Scope:** Paralleles Einspielen mehrerer Sessions; zeitgetreue Antworten
im Replay-Modus.

---

### LH-FA-22 — Wählbares Aufzeichnungsformat

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll Aufzeichnungen wahlweise als Textdatei
(Standard, diff-freundlich und versionierbar) oder als Datenbankdatei (für große
Aufzeichnungen) speichern können. Replay und Einspielen verwenden beide Formate
gleich und erkennen das Format der Datei selbst.

**Akzeptanzkriterien:**

- **Happy Path:** Given der Wunsch nach dem Datenbankdatei-Format, when aufgezeichnet
  wird, then entsteht eine Aufzeichnung als Datenbankdatei, die Replay und Einspielen wie die
  Textaufzeichnung verwenden (siehe Abnahmeszenario 14 in §7).
- **Boundary:** Given kein Wunsch nach einem Format, when aufgezeichnet wird,
  then entsteht eine Textdatei; given eine vorhandene Aufzeichnung, when sie
  verwendet wird, then erkennt das Produkt ihr Format.
- **Negative:** Given eine Datei, die weder eine gültige Text- noch eine gültige
  Datenbankdatei-Aufzeichnung ist, when sie verwendet wird, then wird sie als ungültiges
  Recording erkannt (siehe LH-FA-07).

**Out-of-Scope:** Umwandlung einer Aufzeichnung von einem Format in das andere;
weitere Formate.

---

### LH-FA-23 — Verschlüsselung zum Client

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll im Record- und im Replay-Modus auf Wunsch
verschlüsselte Verbindungen von Clients annehmen, mit einem vom Anwender
bereitgestellten Zertifikat und Schlüssel. Ohne diesen Wunsch lehnt es eine
Verschlüsselungsanfrage wie bisher erkennbar ab und erwartet eine unverschlüsselte
Verbindung (LH-FA-05). Die Verbindung zum PostgreSQL-Server im Record-Modus bleibt
unverschlüsselt. Die Aufzeichnung enthält keine Angaben zur Verschlüsselung der
Client-Verbindung; dieselbe Aufzeichnung ist mit und ohne Verschlüsselung
wiedergebbar.

**Akzeptanzkriterien:**

- **Happy Path:** Given Zertifikat und Schlüssel, when ein Client verschlüsselt
  verbindet, then läuft die Aufzeichnung beziehungsweise Wiedergabe wie bei einer
  unverschlüsselten Verbindung (siehe Abnahmeszenario 15 in §7).
- **Boundary:** Given kein Wunsch nach Verschlüsselung, when ein Client sie
  anfragt, then wird sie erkennbar abgelehnt; given ein Wunsch nach
  Verschlüsselung, when ein Client unverschlüsselt verbindet, then lehnt der
  Recorder die Verbindung ab, sofern der Anwender nicht ausdrücklich beides
  zulässt.
- **Negative:** Given ein fehlendes oder ungültiges Zertifikat oder ein
  nicht passender Schlüssel, when das Produkt startet, then endet der Start mit
  einem eindeutigen Konfigurationsfehler (siehe LH-QA-05).

**Out-of-Scope:** Verschlüsselung zum PostgreSQL-Server im Record-Modus;
Prüfung von Client-Zertifikaten; Zertifikatsverwaltung und -erneuerung.

---

### LH-FA-24 — Vergleich der Antworten beim Einspielen

**Priorität:** SOLL

**Beschreibung:** Das Produkt soll beim Einspielen auf Wunsch die Antworten des
Servers mit der Aufzeichnung vergleichen und Abweichungen melden. Verglichen werden
die Struktur der Antwort einer Interaktion (Art und Reihenfolge der Antworten,
Spaltenbeschreibung einschließlich der Datentypen, Beschreibung der Parameter, der
ausgeführte Befehl, der Transaktionsstatus) und Fehler (Vorhandensein und
Fehlercode). Nicht verglichen werden die Werte und die Anzahl der Zeilen, Hinweise
und Einstellungen des Servers. Mit diesem Wunsch gilt ein Fehler des Servers, den
die Aufzeichnung genauso enthält, als erwartet; der Wunsch ist die ausdrückliche
Wahl nach LH-FA-20. Ohne den Wunsch findet kein Vergleich statt (LH-FA-20).

**Akzeptanzkriterien:**

- **Happy Path:** Given eine Aufzeichnung und ein Server, dessen Antworten in
  Struktur und Fehlern der Aufzeichnung entsprechen, when mit Vergleich
  eingespielt wird, then endet der Lauf mit Erfolg (siehe Abnahmeszenario 16 in
  §7).
- **Boundary:** Given Antworten, die sich nur in Zeilenwerten, Zeilenanzahl oder
  Hinweisen unterscheiden, when verglichen wird, then gilt das nicht als
  Abweichung; given ein Fehler des Servers, den die Aufzeichnung genauso enthält,
  when verglichen wird, then bricht er den Lauf nicht ab und ist keine Abweichung;
  given eine Aufzeichnung, die nicht bis zum Ende der Antwort reicht, when
  verglichen wird, then wird verglichen, soweit sie reicht; given kein Wunsch nach
  Vergleich, when eingespielt wird, then wird nicht verglichen.
- **Negative:** Given eine Antwort mit abweichender Struktur oder einem
  abweichenden Fehler, wozu auch ein Fehler des Servers gehört, den die
  Aufzeichnung nicht enthält, when verglichen wird, then wird die Abweichung mit
  Interaktion und Art der Abweichung gemeldet, und der Lauf endet mit einem
  Fehlerstatus, der sich von dem eines Serverfehlers ohne Vergleich unterscheidet
  (siehe LH-QA-05).

**Out-of-Scope:** Vergleich der Zeilenwerte; Toleranzregeln für einzelne Felder;
Veränderung der Aufzeichnung anhand des Vergleichs.

---

## 4. Nichtfunktionale Anforderungen und Randbedingungen

### LH-QA-01 — Determinismus

- **Anforderung:** Replay-Läufe sollen für identische unterstützte Eingaben und
  identische Aufzeichnungen ein reproduzierbares Verhalten zeigen.
- **Messmethode:** Wiederholter Replay-Lauf mit identischer Aufzeichnung und
  identischer Eingabe, mindestens zehn aufeinanderfolgende Läufe; das
  beobachtbare Verhalten ist in allen Läufen gleich.

### LH-QA-02 — Geringe Eingriffe in die Anwendung

- **Anforderung:** Der Einsatz des Recorders soll keine Änderungen an der
  fachlichen Anwendungslogik erfordern. Die Anwendung soll den Recorder wie einen
  PostgreSQL-Endpunkt ansprechen können.
- **Messmethode:** Wechsel zwischen direkter PostgreSQL-Verbindung und Recorder
  ausschließlich über Verbindungsparameter (siehe Abnahmeszenario 3 in §7).

### LH-QA-03 — Portabilität

- **Anforderung:** Das Werkzeug soll für typische lokale Entwicklungs- und
  CI-Umgebungen bereitstellbar sein.
- **Messmethode:** Lauf der Abnahmeszenarien in einer lokalen
  Entwicklungsumgebung und einer CI-Umgebung; die für die Abnahme verwendeten
  Umgebungen werden bei der Abnahme benannt.

### LH-QA-04 — Automatisierbarkeit

- **Anforderung:** Start, Betrieb und Beendigung müssen für Skripte und
  CI-Systeme geeignet sein.
- **Messmethode:** Vollständig skriptgesteuerter Lauf von Record und Replay
  ohne Eingaben (siehe Abnahmeszenario 5 in §7).

### LH-QA-05 — Nachvollziehbare Fehler

- **Anforderung:** Fehlerzustände müssen so dargestellt werden, dass ein
  Entwickler mindestens folgende Fehlerklassen unterscheiden kann:
  Konfigurationsfehler · Verbindungsfehler · ungültiges oder nicht lesbares
  Recording · nicht unterstützte Protokollinteraktion · Abweichung zwischen
  Replay-Anfrage und Recording.
- **Messmethode:** Je Fehlerklasse ein Szenario, das über Meldung und
  Prozessstatus eindeutig von den übrigen Klassen unterscheidbar ist.

### LH-QA-06 — Wartbarkeit des Recording-Formats

- **Anforderung:** Das persistente Recording muss eine erkennbare Format-
  beziehungsweise Versionsinformation besitzen, sodass zukünftige
  Produktversionen inkompatible Änderungen feststellen können. Die konkrete
  technische Repräsentation ist nicht Bestandteil dieses Lastenhefts.
- **Messmethode:** Jede erzeugte Aufzeichnung trägt die Versionsinformation;
  eine Aufzeichnung mit unbekannter Version wird als solche erkannt.

### LH-RB-01 — Umgang mit sensiblen Daten

- **Art:** rechtlich
- **Vorgabe:** `pgwire-recorder` ist in v1 **nicht dafür verantwortlich,
  sensible Inhalte automatisch zu erkennen oder zu maskieren**. Aufzeichnungen
  können sensible oder vertrauliche Informationen enthalten, etwa SQL-Inhalte,
  zurückgegebene Datensätze, personenbezogene Daten sowie Geschäfts- oder
  Testdaten. Der Anwender ist dafür verantwortlich, geeignete Testdaten zu
  verwenden und Recordings angemessen zu speichern, zu übertragen und zu
  versionieren.
- **Nachweis:** Das Produkt enthält keine Maskierungsfunktion; die
  Verantwortung des Anwenders ist in der Anwenderdokumentation benannt.

### LH-RB-02 — Einsatzumgebung

- **Art:** technisch
- **Vorgabe:** Der primäre Einsatz erfolgt in Entwicklungs-, Test- und
  CI-Umgebungen. Recordings werden unter kontrollierten Testbedingungen
  erzeugt. Die Anwendung kann so konfiguriert werden, dass sie Host und Port des
  Recorders anstelle des PostgreSQL-Servers verwendet. Die für einen
  Replay-Test benötigten Interaktionen wurden zuvor erfolgreich aufgezeichnet.
  Nicht unterstützte PGWire-Funktionen werden mit einem klar erkennbaren Fehler
  abgelehnt.
- **Nachweis:** Abnahmeszenarien in §7 laufen unter diesen Bedingungen.

---

## 5. Globale Out-of-Scope-Punkte

Folgende Funktionen sind ausdrücklich **nicht Bestandteil des verpflichtenden
Funktionsumfangs von v1** und können Gegenstand späterer Versionen werden:

- automatische Maskierung sensibler Daten,
- SQL-semantische Analyse als Voraussetzung für Record oder Replay,
- Veränderung oder Manipulation aufgezeichneter Datenbankantworten,
- Simulation frei definierter PostgreSQL-Antworten ohne vorheriges Recording,
- grafische Benutzeroberfläche,
- Verwaltung eines zentralen Recording-Servers oder Recording-Repositorys.

## 6. Glossar

| Begriff | Bedeutung im Lastenheft |
|---|---|
| Recording / Aufzeichnung | Persistente Datei oder zusammengehöriger Satz von Dateien mit den aufgezeichneten Interaktionen zwischen Client und PostgreSQL |
| Record-Modus | Betriebsart, in der der Recorder zwischen Client und realem PostgreSQL-Server vermittelt und aufzeichnet |
| Replay-Modus | Betriebsart, in der der Recorder Anfragen anhand eines Recordings ohne PostgreSQL-Server beantwortet |
| Einspielen | Betriebsart, in der der Recorder die Client-Anfragen eines Recordings gegen einen PostgreSQL-Server ausführt |
| Simple Query Protocol | PostgreSQL-Protokollvariante für Anfragen als einzelne Textnachricht; Teil von v1 (LH-FA-05) |
| Extended Query Protocol | PostgreSQL-Protokollvariante mit getrennten Parse-/Bind-/Execute-Schritten; Teil von v1 (LH-FA-18) |
| Unterstützt | Teil des in LH-FA-05 und LH-FA-18 beschriebenen Funktionsumfangs von v1 und nicht in §5 ausgeschlossen |

## 7. Abnahmekriterien für v1

Die Produktversion v1 gilt hinsichtlich dieses Lastenhefts als fachlich
abnahmefähig, wenn mindestens folgende Szenarien nachweisbar funktionieren.

### Abnahmeszenario 1 — Erfolgreiches Recording

Eine Testanwendung verbindet sich mit `pgwire-recorder`, der mit einer realen
PostgreSQL-Instanz verbunden ist. Die Anwendung führt mindestens eine
unterstützte Simple-Query-Interaktion erfolgreich aus. Das Tool erzeugt
anschließend eine persistente Aufzeichnung, die für Replay verwendet werden kann.
Bezug: LH-FA-02, LH-FA-06, LH-FA-07.

### Abnahmeszenario 2 — Erfolgreiches Replay ohne PostgreSQL

Die reale PostgreSQL-Instanz ist nicht verfügbar beziehungsweise wird nicht
gestartet. `pgwire-recorder` wird mit der in Abnahmeszenario 1 erzeugten Aufzeichnung im
Replay-Modus gestartet. Die gleiche unterstützte Client-Interaktion kann
erfolgreich durchgeführt werden und liefert das aufgezeichnete, aus Sicht des
Clients relevante Verhalten. Bezug: LH-FA-03, LH-FA-09.

### Abnahmeszenario 3 — Keine Anwendungscode-Integration

Für den Wechsel zwischen direkter PostgreSQL-Verbindung und Recorder ist für den
unterstützten Testfall keine Änderung der fachlichen Programmlogik erforderlich.
Eine Änderung der Verbindungsparameter ist zulässig. Bezug: LH-FA-04, LH-QA-02.

### Abnahmeszenario 4 — Replay-Abweichung

Die Anwendung stellt im Replay-Modus eine unterstützte Anfrage, die nicht zur
verfügbaren Aufzeichnung passt. Der Recorder erkennt die Abweichung und liefert
beziehungsweise signalisiert einen eindeutigen Fehler, statt eine unpassende
aufgezeichnete Antwort zu verwenden. Bezug: LH-FA-10.

### Abnahmeszenario 5 — Automatisierter Betrieb

Record und Replay können vollständig ohne interaktive Eingaben gestartet und in
einem automatisierten Testlauf verwendet werden. Bezug: LH-FA-15, LH-QA-04.

### Abnahmeszenario 6 — Fehlerreplay

Eine während Record auftretende und vom Produkt unterstützte
PostgreSQL-Fehlerantwort wird aufgezeichnet. Beim entsprechenden Replay erhält
der Client ein äquivalentes beobachtbares Fehlerverhalten. Bezug: LH-FA-11.

### Abnahmeszenario 7 — Extended Query

Eine Testanwendung, die das Extended Query Protocol mit Prepared Statements
verwendet (zum Beispiel mit dem Standardmodus eines verbreiteten
PostgreSQL-Treibers), wird über den Recorder gegen eine reale PostgreSQL-Instanz
aufgezeichnet und anschließend ohne PostgreSQL im Replay-Modus ausgeführt; sie
erhält dasselbe beobachtbare Verhalten, ohne dass die Treiberkonfiguration über
Host und Port hinaus geändert wird. Bezug: LH-FA-18, LH-FA-04.

### Abnahmeszenario 8 — Diagnoseausgaben

Es werden drei Störungen herbeigeführt: eine nicht erreichbare PostgreSQL-Instanz
im Record-Modus, eine ungültige Aufzeichnung im Replay-Modus und eine Anfrage,
die nicht zur Aufzeichnung passt. In jedem Fall nennt die Ausgabe des Werkzeugs
die Art und die Ursache des Problems, und die drei Fälle sind an der Meldung
unterscheidbar. Bezug: LH-FA-14.

### Abnahmeszenario 9 — Containerbetrieb

Record und Replay werden mit dem bereitgestellten Docker/OCI-Image in einer
containerisierten Testumgebung gestartet, ohne besondere Sonderkonfiguration
über die Verbindungsparameter hinaus. Der
Replay-Lauf braucht dabei weder eine PostgreSQL-Instanz noch eine Eingabe.
Bezug: LH-FA-16.

### Abnahmeszenario 10 — Maschinenlesbare Konfiguration

Alle Einstellungen, die für einen Record-, einen Replay- und einen
Einspiel-Lauf nötig sind, werden ohne interaktive Eingabe gesetzt, einmal über
Argumente, einmal über Umgebungsvariablen und einmal über eine
Konfigurationsdatei. Wird dieselbe Einstellung über mehrere Wege gesetzt, gilt
die dokumentierte Priorität. Bezug: LH-FA-17.

### Abnahmeszenario 11 — Homebrew

Auf einem Rechner mit Homebrew unter macOS oder Linux wird das Produkt mit den
dokumentierten Befehlen installiert. `pgwire-recorder version` meldet die
installierte Version, und ein Replay-Lauf startet ohne weitere Installation.
Bezug: LH-FA-19.

### Abnahmeszenario 12 — Einspielen

Eine Aufzeichnung mit DDL- und DML-Anweisungen (aus Abnahmeszenario 1) wird
gegen eine leere PostgreSQL-Instanz eingespielt. Danach enthält die Datenbank die
Wirkung der aufgezeichneten Anweisungen, und der Lauf endet mit Erfolg. Bezug:
LH-FA-20.

### Abnahmeszenario 13 — Zeitgetreues Einspielen

Eine Aufzeichnung mit Zeitangaben, in der zwischen zwei Anfragen eine Pause
liegt, wird auf Wunsch zeitgetreu eingespielt. Der Abstand zwischen den beiden
Anfragen ist nicht kürzer als aufgezeichnet. Ohne den Wunsch läuft dieselbe
Aufzeichnung ohne Wartezeit. Bezug: LH-FA-21.

### Abnahmeszenario 14 — Aufzeichnung als Datenbankdatei

Eine Testanwendung wird mit dem Wunsch nach dem Datenbankdatei-Format aufgezeichnet und
anschließend ohne PostgreSQL im Replay-Modus ausgeführt; sie erhält dasselbe
Verhalten wie mit der Textaufzeichnung. Dieselbe Aufzeichnung lässt sich in eine
Datenbank einspielen. Bezug: LH-FA-22.

### Abnahmeszenario 15 — Verschlüsselte Verbindung zum Client

Ein Client verbindet sich verschlüsselt mit dem Recorder im Record-Modus; die
Aufzeichnung entsteht wie bei einer unverschlüsselten Verbindung. Dieselbe
Aufzeichnung wird im Replay-Modus einem verschlüsselt verbindenden Client
geliefert. Ohne Wunsch nach Verschlüsselung wird die Anfrage erkennbar abgelehnt.
Bezug: LH-FA-23.

### Abnahmeszenario 16 — Vergleich beim Einspielen

Eine Aufzeichnung wird mit Vergleich gegen eine PostgreSQL-Instanz eingespielt, die
dieselbe Struktur liefert; der Lauf endet mit Erfolg. Gegen eine Instanz, deren
Antwort in Struktur oder Fehler abweicht, meldet der Lauf die Abweichung und endet
mit einem Fehlerstatus, der sich von dem eines Serverfehlers ohne Vergleich unterscheidet. Bezug: LH-FA-24.

## 8. Historie

Regeln dieser Sektion: Ab Status `Accepted` ist **jede** Änderung an diesem
Dokument eine Vertragsänderung — auch das **Hinzufügen** einer neuen
Anforderung. Sie entsteht **nur** aus einem Change Request, nie aus einem
ADR oder Slice. Fußabdruck pro angenommenem CR: Version-Bump oben plus eine
Zeile hier, mit dem CR unter „Verweis" (Baseline-Regelwerk
`grundlagen-source-precedence.md` §Spec-Stratifizierung).

**Einzige Ausnahme: die Tatsachenberichtigung** — eine Stelle, die etwas anderes
behauptete, als vereinbart war. Sie braucht keinen Change Request, wenn sie hier
**als solche ausgewiesen** ist und keine Aussage einer Anforderung berührt;
„Verweis" trägt dann `—`. Wer eine Sinn-Änderung so etikettiert, hat die
Trennung von Entscheidung und Umsetzung verloren.

Sind Auftraggeber und Entwickler **dieselbe Person**, fehlt nur die
Ticket-Form, nicht der Vorgang. Dann trägt der **Commit** die Trennung: Die
Änderung an dieser Datei liegt in einem eigenen Commit, **vor** dem Slice, der
sie umsetzt — und „Verweis" nennt diesen Vorgang statt eines Tickets.

**Auch hier gilt die Decken-Regel:** keine ADR, kein Slice, kein Carveout,
keine Welle, kein Verweis auf `spezifikation.md` oder `architecture.md` — in
keiner Spalte. Kein Spec-Stratum nimmt seine Historie davon aus. Die Spalte
„Verweis" trägt den **externen** CR; der steht außerhalb des Repos und damit
außerhalb der Referenz-Richtung. Wer im Repo bemerkt hat, dass eine Änderung
nötig wird, hält das auf seiner Seite fest — etwa in der Closure-Notiz des
Slice (Baseline-Regelwerk `modul-03-spec.md`
§Ziel-Form: Akzeptanzkriterium).

| Version | Datum | Änderung | Verweis |
|---|---|---|---|
| 0.1.0 | 2026-10-03 | Initiale Fassung | — |
