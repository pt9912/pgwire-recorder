# Slice slice-v1-abschluss-protokollrand: Protokollversion und Protokollrand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-05.c` · `LH-FA-05.e` · `LH-FA-18.a` · `LH-FA-23.a` · `SPEC-029` · `SPEC-026` · `ARC-006` · `ARC-007`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede Zeile der Tabelle in `LH-FA-05.e` (Protokollversion, `SSLRequest`, `GSSENCRequest`, `CancelRequest`, erste Nachricht und Startnachricht am Rand, unbekannte Sonderanfrage, `COPY`, `FunctionCall`, `NotificationResponse`, `LISTEN`/`NOTIFY`) ist am gebauten Binary durch einen Test belegt, und die Auswahl der Anmelde-Nachrichten, die das Replay emuliert, ist fest und geprüft (`LH-FA-05.b`, *Replay*).

**Stand beim Start (Planner, 2026-10-15, gelesen gegen den Code):** Das Verhalten der meisten Zeilen liefern Vorgänger-Slices und Bestand (`TestSSLUndGSSMitN`, `TestAndereProtokollversion`, `TestCancelRequest`, `TestFremdeErsteNachricht`, `TestVerbindungOhneNachricht`, `TestUnbekannteSonderanfrage`, `TestNichtUnterstuetzteNachricht`, `TestE2ERecordSSLAbgelehnt`, `TestE2ERecordNichtUnterstuetzt`, `TestE2ERecordFremdeAnfrage`, `TestE2ERecordHalbeStartnachricht`; die Anmeldung im Record liefern `slice-v1-abschluss-anmeldung` und `-einspielen-anmeldung`, TLS `slice-v1-abschluss-tls-client`, die Zusage der Versionen `slice-v1-abschluss-postgres-versionen-zusage`). Der Slice ist deshalb vor allem Inventur und Lückenschluss: Der Implementer legt zuerst je Zeile und je Pfad (Happy/Boundary/Negative) und je Modus (`record`, `replay`) die Deckung offen (Abdeckungs-Deklarationen in `docs/user/abdeckung-*.md`) und schließt die Lücken. Vermutete Lücken, ungeprüft: keine Abdeckungs-Deklaration beim Test der unbekannten Sonderanfrage; kein E2E-Test für die andere Protokollversion, `FunctionCall`, `NotificationResponse` und `LISTEN`/`NOTIFY` gegen einen echten Server; die Zeile *Verbindung ohne erste Nachricht* als TCP-Probe im E2E; die Zeilen im Modus `replay`.

**Übernommen aus `slice-v1-abschluss-cancel-ohne-schluessel` und `slice-v1-abschluss-postgres-versionen`** (dort je §1, Abgrenzung; der Geber ist in `slice-v1-abschluss-postgres-versionen-matrix` und `slice-v1-abschluss-postgres-versionen-zusage` aufgegangen, die Abgrenzung steht in `slice-v1-abschluss-postgres-versionen-matrix` §1): der übrige Protokollrand und die Tests der Tabelle in `LH-FA-05.e`; das sind Ziel und DoD-Punkt 1 oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-protokollrand-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Extended-Query-Nachrichten — bereits Gegenstand von welle-extended-query; hier nur der Protokollrand. Die Zeile der Tabelle (*unterstützt*) gilt mit den vorhandenen Tests dieser Welle als belegt und bekommt hier keinen neuen.
- `CancelRequest` mit 12 Byte (ohne Schlüssel) — `slice-v1-abschluss-cancel-ohne-schluessel` (welle-v1-nachtrag; setzt diesen Slice voraus, dort §1 und DoD): Die Zeile der Tabelle ist hier mit der Länge 16 belegt; die Länge 12 lehnt die PGWire-Bibliothek heute ab und ist dort Gegenstand.
- Die TLS-Hälfte der Zeile `SSLRequest` (*mit Zertifikat und Schlüssel angenommen*) und die Anmeldung im Record — Bestand aus `slice-v1-abschluss-tls-client` und `slice-v1-abschluss-anmeldung` (`done/`): Die Tests liegen dort; dieser Slice fügt keinen hinzu und ändert nichts an ihnen.
- Eine Aushandlung der Protokollversion (`NegotiateProtocolVersion`, Version 3.2) und eine Änderung der Liste der PostgreSQL-Hauptversionen — ein anderer Vorgang mit eigener Anforderung bzw. Gegenstand der Matrix (`slice-v1-abschluss-postgres-versionen-matrix`, `done/`): `LH-FA-05.e` lehnt jede andere Startversion als 3.0 mit `PGR-E6002` ab.
- Weitere emulierte Authentifizierungsverfahren im Replay (Passwort, MD5, SCRAM) — wären ein neuer Funktionsumfang statt der Festlegung des vorhandenen; v1 ist keine Sicherheitsgrenze (`LH-FA-05.b`). Wünscht der Nutzer sie, ist das eine Entscheidung vor dem Start (§6).
- Produkt-Code, der kein Verhalten einer Tabellenzeile oder des Replay-Handshakes berührt — Schicht-Abgrenzung: Die Tests dieses Slice sind überwiegend Tests von Bestand; ein Befund außerhalb der Zeilen geht als Befund an den Planner.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Jede Zeile der Tabelle in `LH-FA-05.e` ist durch einen Test belegt (Happy/Boundary/Negative); je Zeile weist die Deckung (`make abdeckung`) den Test nach, und wo eine Zeile in beiden Modi gilt, ist `replay` mit einem Fall belegt oder die Grenze steht in der Spezifikation. Die Zeilen *Extended-Query-Nachrichten* und die TLS-Hälfte von `SSLRequest` zählen mit ihren vorhandenen Tests (§1).
- [ ] Der Replay-Handshake emuliert eine feste Auswahl von Authentifizierungsnachrichten, die die Spezifikation nennt (`LH-FA-05.b`, *Replay*; Entscheidung des Architect vor dem Code, §6), und ein Test belegt sie samt Randformen aus §6; der Handshake für den Go-Client im Standardmodus stammt aus `slice-extended-query-replay`. Die Dokumentation der Auswahl im Handbuch liefert `slice-v1-abschluss-protokollrand-doku`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` | update | Lücken der Tabelle schließen (Tests, im Befund Korrektur), Replay-Handshake festlegen und prüfen |
| `internal/adapters/driven/postgres` | update (nur falls die Inventur eine Lücke zeigt) | `NotificationResponse` als Serverantwort nicht unterstützt (Test in `toResponse` und im Aufbau) |
| `test/integration` | update | E2E-Fälle gegen einen echten Server und das gebaute Binary für die Zeilen ohne E2E; keine Schicht des Hexagons |
| `docs/user/abdeckung-*.md` | update (`make abdeckung`) | erzeugte Tabellen; keine Schicht (`slice-v1-abschluss-protokollrand-doku` §1) |

Die Anmeldung im Record ist Bestand (`slice-v1-abschluss-anmeldung`); der Zeiger „Authentifizierung gegen den Upstream vermitteln“ des Ursprungsplans entfällt.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test spec docs/user`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 1200** geht der Slice zurück, **Meldeschwelle 900** (der Implementer meldet den Stand dem Architect, bevor er weiterarbeitet). Die Marke ist die Schätzung des Planners für rund zehn bis zwölf Tests (Unit und E2E) einschließlich der erzeugten Tabellen in `docs/user`; der Architect gibt seine Schätzung in dieser Metrik vor dem Code an (die Tests tragen rund das Doppelte des Produktcodes) und setzt Marke und Schwelle fest oder hebt sie mit Grund (Vorbild: `slice-v1-abschluss-tls-client` §4). Zweiter Anlass: eine dritte Schicht des Hexagons kommt hinzu (zum Beispiel ein Port oder der Kern), oder ein Schnitt der Zeilen nach Lieferwert (Zeilen des Verbindungsaufbaus; Zeilen der Client-Nachrichten) wird nötig.
- `in-progress` → `open` (blockiert): Eine Zeile der Tabelle zeigt am gebauten Binary ein Verhalten, das die Tabelle anders regelt, und die Spezifikation lässt sich nicht berichtigen, ohne das Lastenheft zu ändern — zuerst die Entscheidung des Nutzers (Carveout nur nach ihr).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — für den Architect, jede vor dem ersten Code-Commit in der Spezifikation (Technik-Stratum, `LH-FA-05.b`/`.e`) zu entscheiden oder als genannt und entschieden zu belegen; der Planner hat sie gegen den Code gelesen, entschieden hat er keine. Eine Entscheidung des Architect an der Spezifikation ist keine Schicht dieses Slice (wie in `slice-v1-abschluss-tls-client` §8).

| Randform | Was offen ist |
|---|---|
| *Auswahl im Replay* | `LH-FA-05.b` sagt nur „in der Implementierung gewählt“. Heute sendet `ARC-006` nach der Startnachricht ohne Anforderung `AuthenticationOk` (`server.go`); die Spezifikation soll das als feste Auswahl nennen (Empfehlung des Planners: nur `AuthenticationOk`, keine Anforderung, kein Passwort). Dazu: Verhalten, wenn ein Client trotzdem eine Nachricht `p` sendet; ein Client mit `require_auth`/Channel Binding (Grenze oder Fall); Benutzer und Datenbank der Startnachricht werden nicht geprüft; gleich für `replay` mit und ohne TLS. |
| *Modus je Zeile* | Die Tabelle nennt keinen Modus. Je Zeile: gilt sie in `record`, in `replay` oder in beiden (Cancel, Version, Sonderanfrage, fremde erste Nachricht, `COPY`, `FunctionCall`, `NotificationResponse`)? Die Spezifikation sagt es, der Test folgt. |
| *Protokollversion* | Werte `0`, `2.0`, `3.1`, `3.2`, `4.0`, `65535.x` (Hauptnummer 1234 ist Sonderanfrage, Zeile darunter); die Fehlerantwort nennt `PGR-E6002` in `record` und `replay`; mit und ohne TLS; Länge der Startnachricht unter 8, genau 8, ohne Parameter, ohne abschließendes Nullbyte; Parameter `_pq_.`-Präfix (Protokollerweiterung) und `options` — ob sie `PGR-E6002` auslösen oder durchgehen. |
| *Sonderanfrage mit unbekanntem Code* | Hauptnummer 1234 mit Nebennummer außer 5679, 5680 und 5678; Länge ≠ 8 bei bekanntem Code ist bei `SSLRequest`/`GSSENCRequest` entschieden (`PGR-E6001`, Länge außerhalb `PGR-W3003`), bei `CancelRequest` mit Länge ≠ 16 nicht (Länge 12 gehört `slice-v1-abschluss-cancel-ohne-schluessel`; andere Längen hier benennen). |
| *`GSSENCRequest`* | Antwort `N` in beiden Modi; danach `StartupMessage`, `SSLRequest` oder weiterer `GSSENCRequest` (mehrfach) — Reihenfolgen und Wirkung mit und ohne TLS; Bytes vor dem `N`. |
| *`COPY`-Nachrichten* | je Richtung: Client sendet `CopyData`, `CopyDone` oder `CopyFail` ohne vorangehende `COPY`-Anfrage (`PGR-E6001`, SQLSTATE `0A000`); Server antwortet auf eine `COPY`-Anfrage mit `CopyInResponse`, `CopyOutResponse` oder `CopyBothResponse`; ob die Verbindung zum Upstream mit `Terminate` geschlossen wird; die Session wird nicht aufgezeichnet (`LH-FA-06` Negative). |
| *`FunctionCall`* | als erste Nachricht nach dem Aufbau, mitten im Aufbau, in einer Extended-Sequenz, im `replay`. |
| *`NotificationResponse`* | wann sie eintrifft: mitten in der Antwort einer Anfrage (heute `PGR-E6001` in `toResponse`), im Leerlauf zwischen zwei Anfragen (der Recorder liest erst bei der nächsten Anfrage), im Aufbau (heute `PGR-E6001`); was der Client erhält; was aufgezeichnet wird. |
| *`LISTEN`/`NOTIFY` als `Query`* | **Möglicher Widerspruch in der Tabelle, vor dem Code zu klären:** Die Zeile sagt „wie jede andere Query aufgezeichnet“, die Zeile darüber „`NotificationResponse` nicht unterstützt“. Ein `NOTIFY` auf einen Kanal, den dieselbe Session mit `LISTEN` hört, liefert PostgreSQL der eigenen Verbindung vor `ReadyForQuery` zu (Vermutung des Planners, mit einer Sonde gegen den gepinnten Server zu belegen): Dann wäre die Folge `LISTEN x; NOTIFY x` in einer Session nicht aufgezeichnet, sondern `PGR-E6001`. Der Architect legt die Fälle fest (`LISTEN` allein, `NOTIFY` allein, beides in einer Session, zwei Sessions) und berichtigt die Zeile im Wortlaut (`BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut`, 1×). Das Verhalten hängt nicht von der Serverversion ab (zu belegen gegen 14 und 18, `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×). |
| *Verbindung ohne erste Nachricht / Halbe Startnachricht* | TCP-Probe (Verbindung und sofortiges Schließen, auch mit Reset), 1 bis 7 Bytes, genau 8 Bytes mit Länge über 8; in beiden Modi; mit `--record-empty-sessions` (entschieden in `slice-v1-abschluss-sessions`, hier nur geprüft). |
| *Aufrufart: gleichzeitig* | Eine Zeile gilt unabhängig von anderen Verbindungen; ein Fall mit zwei gleichzeitigen Verbindungen, eine mit Protokollrand und eine regulär — die reguläre bleibt unberührt (keine Zusage über Ordnung zwischen Verbindungen). |

Zu jeder Zeile der Tabelle und zu jeder Randform gilt `AGENTS.md` §3.10: eine Mutation, die nur sie bricht, vom Implementer selbst rot gesehen; Aufzählungen (Versionswerte, Längen, Nachrichtenarten) zählen je Element.

**Gebundene Sendungen:** Keine eingegangen (`grep -rn "slice-v1-abschluss-protokollrand" docs/plan/planning/next docs/plan/planning/open` findet nur `-protokollrand-doku` (Folge), `-cancel-ohne-schluessel` und `-verbindungsende-format` (Reihenfolge, kein Gegenstand)). Zeigt sich im Lauf eine Sendung an `slice-v1-abschluss-cancel-ohne-schluessel`, geht sie an den Planner.

**Risiken:**

- Die Inventur findet mehr Lücken als der Plan vermutet, und die Zahl der Tests sprengt Größe oder Review-Sitzung (`BEO-REPO/slice-waechst-durch-uebernahmen`, `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, verkörpert) — **Ausgang:** offen bis Closure; Gegenmittel: Marke und Meldeschwelle aus §4, Schnitt nach Lieferwert (Zeilen des Aufbaus; Zeilen der Client-Nachrichten).
- `LISTEN`/`NOTIFY` in einer Session ist nach der Sonde nicht aufgezeichnet, und die Zeile der Tabelle muss berichtigt werden — **Ausgang:** offen bis Closure; Gegenmittel: Berichtigung der Zeile als Entscheidung des Architect vor dem Code (§6 oben), kein Eingriff in das Lastenheft nötig; wünscht der Nutzer, dass `NotificationResponse` aufgezeichnet wird, ist das eine neue Anforderung und eine Entscheidung des Nutzers.
- Die Auswahl im Replay genügt einem Treiber nicht, der Anmeldung verlangt (`require_auth`) — **Ausgang:** offen bis Closure; Gegenmittel: Grenze in der Spezifikation nennen (v1 ist keine Sicherheitsgrenze), Fall im Test.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-protokollrand-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 2. Schichten: zwei Schichten (PGWire-Adapter, Upstream-Adapter).

Nachgezählt beim Start (Planner, 2026-10-15, gegen den heutigen Code und `SPEC-029` nach den Matrix- und Zusage-Slices, `AGENTS.md` §3.9, §3.13): Liefer-Punkte 2 (Tabelle `LH-FA-05.e` getestet; Replay-Handshake festgelegt und geprüft); die Randformen aus §6 gehören zu ihnen und fügen keinen hinzu. Schichten: zwei, PGWire-Adapter (`ARC-006`) und Upstream-Adapter (`ARC-007`, nur bei einer Lücke der Inventur); `test/integration` ist Test des gebauten Binaries und keine Schicht des Hexagons, die erzeugten Tabellen in `docs/user/abdeckung-*.md` zählen nicht (`slice-v1-abschluss-protokollrand-doku` §1), Meldungscode-Konstanten zählen nicht, die Festlegung in der Spezifikation durch den Architect vor dem Code nicht (§6). Berührt der Implementer Kern, Ports oder CLI, sind es drei: Rückführung nach §4. `SPEC-029` wurde durch `slice-v1-abschluss-postgres-versionen-zusage` geändert (Liste 14 bis 18, Referenz 17); der Gegenstand dieses Slice (Tabelle `LH-FA-05.e`) ist davon nicht betroffen, §1 nennt die Versionen nicht mehr als Ziel.

**Vorgelagert — Sichtung für die Berührung (Planner, 2026-10-15):** Das Register trägt inzwischen `BEO-REPO/` (die Zeile oben galt zum Datum des Plans). Treffer: `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` (1×, offen) — die Zeilen `NotificationResponse` und `LISTEN`/`NOTIFY` (§6); `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×, offen) — die Sonde und E2E-Fälle gegen einen echten Server (§6); `BEO-REPO/spec-randform-erst-im-review-entschieden` (29×), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (35×), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (16×), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (56×), `BEO-REPO/plan-folgt-korrektur-nicht` (47×) und `BEO-REPO/slice-waechst-durch-uebernahmen` (5×) — verkörpert, die Regeln in `AGENTS.md` §3.9 bis §3.13 tragen; `BEO-REPO/verhalten-nur-unter-linux-geprueft` (3×, offen, Lücke): berührt nicht, der Slice prüft TCP-Verhalten, keine Eigenschaft des Dateisystems. Keiner der nicht verkörperten Einträge erreicht mit diesem Slice 3× neu; erreicht `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` mit der Berichtigung der Zeile 2×, bleibt er unter der Schwelle.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
