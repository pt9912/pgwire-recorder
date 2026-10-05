# Validierung: slice-extended-query-record, Herunterfahren im Record — 2026-10-05

**Rolle:** Validator (Modul 8). Hier geht es darum, ob das Richtige gebaut wurde, gemessen am Bedarf. Ob richtig gebaut wurde, hat die [Verifikation](2026-10-05-verifikation-slice-extended-query-record.md) geprüft; ihr Ergebnis ist hier Eingang, kein Beleg.

**Gegenstand:** wie sich der Record-Modus beim Herunterfahren verhält, Stand HEAD `22b535a`. Die Fragen kommen aus Plan §6 von `docs/plan/planning/in-progress/slice-extended-query-record.md`: F-309 und F-317 (Herunterfahren ohne Obergrenze, mit dem Material aus V-24) sowie F-318 (Ende einer Session bei einem Client, der nicht liest). Den übrigen Lieferwert des Slice (Extended Query im Record) bewertet dieser Beleg nicht.

**Bedarfsquellen:**

- Lastenheft §1 (Ausgangssituation, Anwendungsfälle) und §2 (Stakeholder: Teams mit CI-Umgebungen, Maintainer von Testinfrastrukturen)
- [LH-FA-02](../../spec/lastenheft.md#lh-fa-02--record-modus) Boundary („der Recorder bleibt nicht hängen“)
- [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary und Negative
- [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben)
- [LH-FA-15](../../spec/lastenheft.md#lh-fa-15--ci-eignung)
- [LH-FA-16](../../spec/lastenheft.md#lh-fa-16--container-eignung)
- [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)
- [LH-QA-04](../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)
- [Abnahmeszenario 5](../../spec/lastenheft.md#abnahmeszenario-5--automatisierter-betrieb) und [Abnahmeszenario 9](../../spec/lastenheft.md#abnahmeszenario-9--containerbetrieb)
- aus der Spezifikation `LH-FA-13.a`, `LH-FA-13.b`, `LH-FA-02.b`, `LH-FA-07.a` und `LH-FA-18.a`
- das [Benutzerhandbuch](../user/benutzerhandbuch.md) und die [README](../../README.md)

Den realen Einsatz hat der Auftrag geschildert: Ein anderes Repo setzt den Recorder in CI- und E2E-Setups zwischen Anwendung und PostgreSQL ein. Gestartet wird er dort als Prozess oder Container, beendet per `SIGTERM`, und der Container hat eine Stopp-Frist, etwa 10 s bei Docker bis zum `SIGKILL`.

**Grenze dieses Belegs:** Er misst gegen den dokumentierten Bedarf und gegen die Schilderung im Auftrag. Eine Aussage des Auftraggebers liegt nicht vor. Wo dieser Beleg eine Lastenheft-Änderung empfiehlt, entscheidet der Auftraggeber. Nach Modul 8 liegt ein Validierungsbeleg gewöhnlich außerhalb des Repos; dass er hier unter `docs/reviews/` steht, hat der Auftrag so festgelegt.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

---

## 1. Beobachtetes Verhalten

**Aus dem Code** (gelesen, nicht verändert):

- Auf das erste Signal schließt `bootstrap.record` den Listener. Danach wartet es, bis alle Sessions beendet sind, und schreibt anschließend (`Finish`). Dieses Warten hat keine Frist.
- Jede Session schreibt die Aufzeichnung bei ihrem eigenen Ende, als Ganzes und atomar (`CloseSession`, `LH-FA-07.a`).
- Eine Session ohne laufende Interaktion endet sofort beim Herunterfahren. Eine Session mit laufender Interaktion endet erst nach deren `ReadyForQuery`. Das gilt für eine einfache Anfrage, bei der `Query` synchron auf den Upstream wartet, genauso wie für eine Extended-Interaktion. Fehlt das `Sync`, etwa nach einem `Flush`, oder antwortet der Upstream nicht, wartet die Session beliebig lange.
- Ein zweites Signal beendet den Prozess sofort (`main.go`, Standardverhalten von Go). Für `record` steht das weder in der Spezifikation noch im Handbuch; die Spezifikation sagt es nur für `play` (`LH-FA-20.a` Schritt 7).
- Beim Beginn des Herunterfahrens gibt der Recorder keine Log-Zeile aus, auch keine über wartende Sessions.

**Sonde** (Kopie per `git archive HEAD`, Stufe `runtime` als eigenes Image, eigenes Docker-Netz, PostgreSQL aus dem gepinnten Image von `harness/mk/integration.mk`, psql als Client; danach alles gelöscht, Arbeitsbaum unverändert):

| Schritt | Beobachtung |
|---|---|
| Session A: `select 1`, danach Leerlauf | — |
| Session B: `select 2`, danach `select pg_sleep(40)` (läuft beim Stopp) | — |
| `docker stop -t 10` auf den Recorder | Der Stopp dauerte 10 s, dann `SIGKILL`, Exit-Code `137` |
| Log auf `stderr` | nur `record gestartet`; keine Zeile zum Herunterfahren, zum Warten oder zum Verlust |
| Aufzeichnung | enthält Session A. **Session B fehlt ganz**, auch ihre abgeschlossene Interaktion `select 2` |

Die Sonde nimmt den Weg über eine einfache Anfrage, weil psql keine `Flush`-Gruppe senden kann. Laut Code verhält sich eine Extended-Interaktion ohne `Sync` genauso. Die Ursache liegt also schon vor diesem Slice, im Pfad der einfachen Anfrage. Der Slice fügt nur weitere Auslöser hinzu: `Flush` ohne `Sync` (F-309) und `Send` an einen Server, der nicht liest (F-317).

---

## 2. Urteile

### Frage 1 — F-309 und F-317: unbegrenztes Warten beim Herunterfahren

**Urteil: trifft den Bedarf nicht.** Hier ist die Verifikation grün und die Validierung rot: Der Plan ist erfüllt, aber der Bedarf ist es nicht.

Begründung am Bedarf:

- Laut [LH-QA-04](../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) muss die *Beendigung* für Skripte und CI geeignet sein. Nach [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary liefert ein kontrolliertes Beenden einen Prozessstatus, der Erfolg oder Fehler sagt. Im realen Einsatz wird ein Beenden ohne Obergrenze aber zu einem unkontrollierten: Nach der Stopp-Frist kommt `SIGKILL`, und laut `LH-FA-13.b` gilt dann kein Exit-Code der Spezifikation.
- Die Folge ist schlimmer als ein Hänger: Die Aufzeichnung bleibt gültig, ihr fehlt aber still eine ganze Session samt deren abgeschlossenen Interaktionen. Kein Log sagt das, und `137` unterscheidet sich nicht von einem Absturz. Mit [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben) ist das nicht vereinbar. Ein CI-Lauf, der das Ergebnis eines `docker compose down` nicht auswertet, versioniert so eine unvollständige Aufzeichnung als Testgrundlage.
- [LH-FA-02](../../spec/lastenheft.md#lh-fa-02--record-modus) Boundary verlangt ausdrücklich, dass der Recorder nicht hängen bleibt. Gemeint ist dort die Verbindung ohne Anfrage, die Absicht ist aber allgemein.
- Gegenargument, geprüft: Ein Test, der vor dem `SIGTERM` zur Ruhe kommt, ist nicht betroffen, und ein Client, der nach `Flush` ruht, ist bei verbreiteten Treibern selten. Der Bedarf entscheidet sich aber am Fehlerfall: Ein Test bricht ab, eine Anwendung hängt in einer Sperrwartung gegen eine Verbindung außerhalb des Recorders, oder das Netz zum Upstream reißt ab. Gerade dann braucht CI ein deutliches Ergebnis statt eines stillen Verlusts.

**Empfehlung:**

1. **Obergrenze:** eine Frist für das Herunterfahren, gezählt ab dem Signal (nicht ab dem ersten `Shutdown` einer Session, siehe Frage 2). Nach Ablauf endet jede noch laufende Session zwangsweise, und zwar so wie bei einem Abbruch nach `LH-FA-02.b`: Die laufende Interaktion wird nicht übernommen, die abgeschlossenen Interaktionen der Session bleiben erhalten. Danach wird die Aufzeichnung geschrieben.
2. **Code und Exit:** Ein Verbindungsfehler der Klasse 4 wird gemerkt, Exit-Code `4`. Der Meldungscode sollte ein eigener sein (Vergabe nach `SPEC-034`), nicht `PGR-E4003`: „Verbindung unerwartet beendet“ führt in die Irre, wenn der Recorder selbst beendet hat. Je abgebrochener Session nennt das Log die Session und die verworfene Interaktion.
3. **Größenordnung, einstellbar:** Default 5 s. Das liegt sicher unter den 10 s von Docker und lässt Zeit für Schreiben und Exit. Einstellbar per Option und Umgebungsvariable, weil [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) verlangt, dass jede Einstellung für Testläufe nicht-interaktiv setzbar ist. Wer beim Herunterfahren lange Anfragen aufzeichnet (Migrationen, Lasttests), braucht mehr Zeit; `0` als „ohne Frist“ ist vertretbar. Eine feste Frist trifft den Bedarf nicht, weil Kubernetes (30 s) und Docker (10 s) unterschiedliche Stopp-Fristen setzen.
4. **Zweites Signal:** Es sollte die Frist sofort ablaufen lassen, also zwangsweise beenden, schreiben und mit Exit-Code enden, statt den Prozess hart zu beenden. Ein harter Abbruch verliert genau die Sessions, auf die gewartet wird. Mindestens gehört das heutige Verhalten in Spezifikation und Handbuch.
5. **Diagnose:** Beim Beginn des Herunterfahrens eine Info-Zeile mit der Zahl der Sessions, auf die gewartet wird.

**Zielort:**

- **Lastenheft:** eine Ergänzung an [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary, etwa: „Das kontrollierte Beenden endet in einer begrenzten, einstellbaren Zeit; was bis dahin nicht abgeschlossen ist, wird nicht übernommen und als Fehler signalisiert.“ Das Lastenheft steht auf `Draft`, die Änderung braucht also keinen Change Request. Weil sie aber eine zugesagte, messbare Eigenschaft betrifft, sollte der Auftraggeber sie bestätigen (Spec-Validierung vor der Umsetzung).
- **Spezifikation:** `LH-FA-13.a` und `LH-FA-13.b` (Frist, Zwangsende, Meldungscode, zweites Signal), dazu die Optionstabelle.
- **Umsetzung:** `slice-v1-abschluss-betrieb`. Sein Ziel ist genau das kontrollierte Herunterfahren mit Exit-Code und einstellbaren Optionen, und die Ursache betrifft beide Interaktionsarten, nicht nur Extended Query. Der Slice nennt die Frist heute nicht. Sein DoD-Punkt 1 („Herunterfahren schreibt das Recording atomar und liefert den Exit-Code …“) kann sie aufnehmen. Wird der Punkt dadurch zu groß, gehört die Frist in einen eigenen Slice der welle-v1-abschluss.
- **Dieser Slice:** nicht die Umsetzung, aber zwei Dinge vor der Closure. Erstens das Risiko F-309 mit dem Ausgang „weiter offen“ an `slice-v1-abschluss-betrieb` binden, und dieser Slice nimmt es in §1, §6 und DoD auf. Zweitens die Hinweise im Handbuch (§4, Hinweise zum erweiterten Protokoll und zur Testautomatisierung) für die Zwischenzeit ergänzen: Die Stopp-Frist des Containers muss länger sein als die längste Anfrage, die beim Beenden noch laufen kann. Bei einem `SIGKILL` fehlt jede noch wartende Session ganz, während bereits beendete Sessions erhalten bleiben. Ein zweites Signal beendet sofort. Das Compose-Beispiel kann `stop_grace_period` zeigen.

### Frage 2 — V-24: Beginn des Herunterfahrens

**Urteil: für Nutzer akzeptabel. In der Spezifikation ist es mit inneren Zeitpunkten formuliert; das sollte sich ändern.**

- Für einen Nutzer ist der Signalzeitpunkt gegenüber dem Datenverkehr ohnehin nicht genau festzulegen. Ob eine Nachricht im Fenster zwischen Signal und Wecken noch verarbeitet oder verworfen wird, kann er nicht beobachten und braucht er auch nicht. Was er braucht, ist erfüllt: Die Aufzeichnung enthält nur abgeschlossene Interaktionen, eine verworfene Nachricht wird nicht halb aufgezeichnet, der Exit-Code bleibt `0`, und der Client sieht das Verbindungsende. Das Handbuch sagt das richtig: „eine neue Anfrage oder Folge beginnt danach nicht mehr, und die Verbindung wird geschlossen“. Für einen Test, der vor dem Signal zur Ruhe kommt, hat das Verhalten keine Folgen.
- Nicht akzeptabel ist nur die Kopplung an Frage 1. Steht eine Session in `Send` an einem Server, der nicht liest (F-317), beginnt für sie das Herunterfahren nie. Das löst die Frist aus Frage 1, wenn sie ab dem Signal zählt.
- Die Spezifikation sagt mit „schon vollständig gelesen“ (`LH-FA-13.a`) mehr zu, als ein Nutzer prüfen kann, und hängt das Verhalten am Lesepuffer auf.

**Empfehlung:** `LH-FA-13.a` so umformulieren, dass beobachtbar bleibt, was gilt. Etwa: „Eine Client-Nachricht, die um das Signal herum eintrifft, wird entweder noch verarbeitet oder nicht weitergeleitet; danach beginnt keine neue Interaktion.“ Wie die Absprache mit der Frist aussieht, bleibt Teil von Frage 1. Optional, ohne dass der Bedarf es verlangt: dem Client vor dem Schließen die Fehlerantwort schicken, die PostgreSQL beim Herunterfahren sendet (`FATAL`, SQLSTATE `57P01`). **Zielort:** die Spezifikation in diesem Slice, weil er den Satz eingeführt hat und `AGENTS.md` §3.9 den Plan der Korrektur folgen lässt. Am Lastenheft ist nichts zu ändern.

### Frage 3 — F-318: Verbindung schließen, Fehlerantwort höchstens 1 s

**Urteil: trifft den Bedarf.**

- Ein Client, der nicht liest, hält die Session nicht mehr auf. Damit wird [LH-FA-02](../../spec/lastenheft.md#lh-fa-02--record-modus) Boundary (nicht hängen) erfüllt, und der Prozess kann sich beenden ([LH-QA-04](../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)).
- Für CI zählen Exit-Code und Log, nicht die Antwort beim Client. Der Recorder merkt sich den Fehler und protokolliert ihn, *bevor* er die Fehlerantwort schreibt (`fail` → `note`). Fällt die Antwort bei einem Client aus, der nicht liest, geht deshalb weder die Klasse noch der Exit-Code verloren ([LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Negative).
- Eine Sekunde reicht für eine einzelne kurze Nachricht an einen lesenden Client bei weitem. Sie greift nur, wenn der Puffer zum Client voll ist, und dann hilft auch eine längere Frist nicht. Weil die Sessions parallel enden, addiert sich die Frist nicht. Einstellbar muss sie nicht sein.
- Dass die Fehlerantwort hinter einem noch laufenden Stapel stehen kann (Plan §6), betrifft nur einen Client, der ohnehin abgebrochen wird. Am Bedarf ändert das nichts.

**Empfehlung:** keine Änderung. Für Frage 1 gilt: Die Sekunde ist Teil des Budgets innerhalb der Frist für das Herunterfahren. **Zielort:** keiner. Das Risiko kann bei der Closure dieses Slice den Ausgang „entfallen“ (vom Bedarf getragen) bekommen.

---

## 3. Nebenbefund

Die [README](../../README.md), Abschnitt „Was kann ich heute tun?“, sagt noch, das erweiterte Protokoll fehle. Für den Record-Modus gilt das nach diesem Slice nicht mehr. Wer vom anderen Repo aus nur die README liest, unterschätzt den Stand. **Zielort:** dieser Slice, bei der Closure.

## 4. Übergabe an den Planner

| Nr. | Gegenstand | Urteil | Zielort |
|---|---|---|---|
| 1 | F-309/F-317: kein Ende des Herunterfahrens in Sicht | rot | Lastenheft [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary (Draft, Bestätigung durch den Auftraggeber), dann `LH-FA-13.a`/`.b`, Umsetzung in `slice-v1-abschluss-betrieb` (DoD-Punkt 1 erweitern oder eigener Slice); in diesem Slice §6 binden und Hinweise im Handbuch |
| 2 | V-24: Beginn des Herunterfahrens | grün, Formulierung der Spezifikation anpassen | `LH-FA-13.a` in diesem Slice |
| 3 | F-318: Schließen, Fehlerantwort 1 s | grün | keiner; Ausgang „entfallen“ |
| — | README-Stand | Drift | dieser Slice, Closure |

Die Closure dieses Slice blockiert keiner dieser Punkte, solange Punkt 1 als gebundenes Risiko mit Folge-Slice weitergeht und das Handbuch das Verhalten in der Zwischenzeit benennt.
