# Verifikation slice-v1-abschluss-protokollrand-doku

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `cf08408` (Basis) bis `ffcfdbf` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-protokollrand-doku` (Plan unter `docs/plan/planning/in-progress/`) gegen DoD (§2), Plan §1 (Sendung des Gebers `slice-v1-abschluss-protokollrand`), §3, §6 (R1 bis R8, P1 bis P7), §7, §8, `spec/spezifikation.md` (`LH-FA-05.e`, `LH-FA-05.b`, `LH-FA-23.a`, *Zählung der Verbindungen*) und `AGENTS.md` §3.9, §3.11, §3.13, §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-protokollrand-doku.md`; gegen ihn ist nicht verifiziert, F-1 und F-2 sind unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-protokollrand`. Belege sind in diesem Lauf selbst erzeugt: Binary aus `make build` (Image `pgwire-recorder:dev`, Image-ID `9a27f080a356`, wie in §7), mit `docker cp` aus dem Image in das Scratchpad kopiert und dort als Prozess gegen `postgres:17-alpine@sha256:b0f9560a…` in einem eigenen Docker-Netz gefahren (Port auf `127.0.0.1`); Clients: ein Python-Skript mit rohen Nachrichten, `psql` 17.11 und `psql` 18.6 (`--network host`). Ein Lauf zur Zeit, kein `pkill`, Prozesse per PID beendet. Danach 0 Container, nur die drei Standard-Netze; `git status` leer.

**Verdikt: PASS mit Befunden.** Alle fünf Inhalte der Sendung stehen im Handbuch und im README, und jede Verhaltensaussage, die ich gegen das gebaute Binary gefahren habe, stimmt. 0 HIGH, 0 MEDIUM, 3 LOW, 3 INFO. Die drei LOW-Befunde sind Lücken oder zu weite Sätze neben dem Neuen, kein falscher Satz; sie gehören vor die Closure zum Implementer oder werden als Grenze ins Register getragen.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren oder gelesen) |
|---|---|---|
| 1: Handbuch beschreibt die fünf Inhalte im Ist-Zustand; Optionen gegen `--help`, Codes im Katalog, keine Verweise, kein „noch nicht/kommt/geplant" | bestätigt, mit V-1 und V-2 | Fünf Inhalte siehe nächster Abschnitt. Neu geschriebene Zeilen (`git diff -U0 cf08408..HEAD -- docs/user README.md`, nur `+`-Zeilen) auf `spezifikation\|lastenheft\|ADR-\|slice\|welle\|review\|LH-\|SPEC-\|noch nicht\|geplant\|künftig\|bisher\|kommt\|später\|derzeit\|aktuell`: 0 Treffer. Das ganze Handbuch zeigt dieselben Treffer wie vorher (Zeilen 175, 825, 1236, 1341, 1346, 1504, 1726: Zieldatei „existiert noch nicht", „noch nicht gültig", „kommt" in Sachaussagen; keiner neu). Keine neue Option, kein neuer Codeblock (Diff gelesen). Meldungscode `PGR-E6001` bis `PGR-E6003`, `PGR-W3001` stehen im Katalog (`make gates-schnell` fährt `docs-check`, siehe unten). |
| 2: README nennt dasselbe in Kürze, Verweise ohne Stand | bestätigt | Absatz „Am Rand des Protokolls" und der geänderte Satz zu `COPY`/Benachrichtigung gelesen; keine Chronik, kein Zielstand, Verweis auf das Benutzerhandbuch. |
| 3: `make gates` grün | durch den Hook getragen (nicht selbst voll gefahren) | Der Diff `cf08408..HEAD` berührt nur `README.md` (Wurzel), `docs/user/benutzerhandbuch.md`, `docs/plan/planning/` und `docs/reviews/`: vier Pfade der Positivliste (`harness/README.md`, Zeile `make stop-hook-gegenprobe`). Mit meinem Report kommt nur eine weitere Datei unter `docs/reviews/` hinzu. Ein voller `make gates` ist deshalb nicht nötig, sofern die Dateiliste des vollen Laufs (`.harness/state/gates-passed.dateiliste`, 08:30) sich nur auf diesen Pfaden unterscheidet; das entscheidet der Stop-Hook, ich habe die Liste nicht selbst verglichen. `make gates-schnell` einmal gefahren: Exit 0; `baseline-verify` v6.18.0 OK (54 Dateien), `d-check` 788 Dateien, 0 Befunde, `gesamt: 0 Befund(e)`; Stand `.harness/state/gates-schnell-passed.diffsha` geschrieben (ignoriert, `git status` leer). |
| 4: Review durchgeführt | bestätigt | Report liegt vor (0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO); F-1 und F-2 unten nachgeprüft. |
| 5 bis 8: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt Belege des Implementers, aber noch keine Closure-Notiz; §2 ohne Häkchen; §6-Risiken tragen „offen bis Closure". Handlung des Planners, Kandidaten siehe unten. |

## Die fünf Inhalte der Sendung gegen das gebaute Binary

| Inhalt | Im Handbuch / README | Mein Lauf | Urteil |
|---|---|---|---|
| (1) `_pq_.` ist `PGR-E6001` in `record` und `replay`, auch mit TLS-Pflicht statt `PGR-E6003` | Hinweise zum Aufzeichnen, TLS-Abschnitt, Tabelle `PGR-E6001`/`PGR-E6003`, README | `record --tls-cert` (Upstream `127.0.0.1:1`, nicht erreichbar): `_pq_.x=1` unverschlüsselt → `PGR-E6001`, ohne `_pq_.` → `PGR-E6003`, Exit 6, kein `PGR-E4002` (kein Verbindungsaufbau zur Datenbank). `replay --tls-cert`, beide Zuordnungen: gleiche Antworten. Ohne TLS in beiden Modi `PGR-E6001`. | bestätigt |
| (2) Zählung abgelehnter Verbindungen; im `replay` keine verbrauchte Session | Hinweise zum Aufzeichnen (Aufzählung), Hinweise zum Wiedergeben (Zählung), README | `record --record-empty-sessions`: Version 3.2, `_pq_.x`, Rumpf Länge 9, danach eine reguläre Verbindung: Datei mit `empty_sessions: true` und genau 1 Sitzung, Exit 6. `replay`, `first-request` und `connection`, TLS-Pflicht: Klartext (`PGR-E6003`) und Klartext mit `_pq_.` (`PGR-E6001`) verbrauchen keine Session, die beiden regulären TLS-Verbindungen danach erhalten `SELECT 1;` und `SELECT 2;`, Exit 6. | bestätigt; zur Vollständigkeit bei `connection` V-1 |
| (3) `NotificationResponse`, auch durch `LISTEN`/`NOTIFY`, ist `PGR-E6001`, Session nicht aufgezeichnet | Hinweise zum Aufzeichnen, Treiber-Hinweise, Tabelle | `LISTEN k; NOTIFY k;` in einer Anfrage, `LISTEN k;` dann `NOTIFY k;` in zwei Anfragen: Fehlerantwort `PGR-E6001` an die Anfrage, Verbindung nicht in der Datei, Exit 6. Zwei Verbindungen: die sendende (`NOTIFY k;`, `SELECT 13;`) in der Datei, die zuhörende bekommt `PGR-E6001` erst auf die nächste Anfrage (`SELECT 14;`) und fehlt. `LISTEN k;` allein: aufgezeichnet, Exit 0. Benachrichtigung nach der letzten Anfrage der zuhörenden Verbindung: beide Sitzungen in der Datei, Exit 0. | bestätigt; zu weite Sätze V-3 |
| (4) Feste Auswahl im Replay: `AuthenticationOk` ohne Anforderung, `p` nach dem Aufbau `PGR-E6001`, `require_auth` als Grenze | Hinweise zum Wiedergeben, Passwort-Anmeldung, Datenschutz-Abschnitt | `replay` mit und ohne TLS, `psql` 17.11, Benutzer `gibtesnicht`, Datenbank `nix`: `require_auth=none` → Ergebnis `1`, Exit 0; `password`, `md5`, `scram-sha-256` → Abbruch im Client („server did not complete authentication"), `replay` Exit 0, keine Meldung zum Abbruch (die einzige `PGR-`-Zeile ist `PGR-W2001`, eine nie zugeordnete zweite Session der Zwei-Session-Aufzeichnung, nicht der Abbruch). `p` und `CopyData` nach dem Aufbau mit unbekanntem Benutzer, Datenbank und `options`: `PGR-E6001`, Verbindung endet, Exit 6, in `record` und `replay`. | bestätigt |
| (5) COPY-Antworten des Servers | Hinweise zum Aufzeichnen, Tabelle | `COPY t FROM STDIN;` und `COPY (SELECT 1) TO STDOUT;`: Fehlerantwort `PGR-E6001` (`CopyInResponse` bzw. `CopyOutResponse`) statt der Antwort, nur die dritte Verbindung (`SELECT 17;`) in der Datei, Exit 6. | bestätigt |

Dazu die Randformen, die der Auftrag nannte:

| Punkt | Ergebnis |
|---|---|
| `CopyDone`/`CopyFail` von der Anwendung | Binary: `PGR-E6001` (*pgproto3.CopyDone / CopyFail*). Im Handbuch nur `CopyData` genannt, `grep -n "CopyDone\|CopyFail"` leer: wie R6 verlangt, kein falscher Satz. |
| `NotificationResponse` im Aufbau | `grep` im Handbuch: kein Satz dazu (nur „sendet einer Verbindung eine Benachrichtigung"), wie R7. |
| `CancelRequest` mit Länge außer 16 | Nicht im Handbuch (`grep -n "Länge 12\|12 Byte"` leer); die Zeile `PGR-W3001` ist unverändert. Länge 12 steht in `slice-v1-abschluss-cancel-ohne-schluessel` §1, DoD und §3. |
| `--record-empty-sessions` mit abgelehnten Verbindungen | siehe (2): genau 1 Sitzung; ohne Verallgemeinerung auf die Verbindung mit `p`/`CopyData`: V-2. |
| `require_auth=none` | Gelingt, Ergebnis `1`, Exit 0 (siehe (4)); das Handbuch nennt nur `password`, `md5`, `scram-sha-256`. |
| P6 (`CancelRequest` mit Schlüssel einer anderen Session) | Nicht baubar, wie §7 sagt: die Anwendung erhält `R`, 14 mal `S`, `Z`, kein `K` (`BackendKeyData`). Der Halbsatz „wertet weder Prozess-ID noch Schlüssel aus" und der Satz zur anderen Verbindung stehen nicht im Handbuch (`grep -n "Prozess-ID"` leer). |
| P7 | `psql` 18.6: `max_protocol_version=3.0` verbindet mit `record` und `replay` (Ergebnis `1`); `latest` und `3.2` → `PGR-E6002` („Protokollversion 3.2"), beide Modi Exit 6. Der Hinweis im Handbuch (`PGR-E6002`, Treiber-Hinweise) ist belegt. |

## Nachprüfung von F-1 und F-2 des Reviews (unabhängig)

| Finding | Mein Befund |
|---|---|
| F-1 (LOW, „jede Verbindung" mit `AuthenticationOk` gegen die abgelehnten Verbindungen derselben Seite) | Beobachtet: eine abgelehnte Verbindung (3.2, `_pq_.`, Rumpf Länge 9, Klartext bei TLS-Pflicht) erhält als erste Nachricht `E`, nicht `R0`; eine angenommene beginnt mit `R0`. Der Satz steht jetzt mit „jede angenommene Verbindung" in Handbuch und README und nennt die sofortige Fehlerantwort. Die zweite Stelle gleichen Merkmals (§6 „bedient eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank") ist gleich gefasst. Behoben. Die Klasse (Quantor gegen Ausnahme derselben Seite) kehrt in V-2 und V-3 wieder. |
| F-2 (INFO, „im Container") | Der Satz lautet „unter Linux geprüft; für andere Systeme ist es nicht geprüft", ohne „im Container". Siehe nächster Abschnitt. |

## Der Linux-Satz (Auftrag 3)

Handbuch (Hinweise zum Aufzeichnen): „Das Verhalten am Protokollrand ist unter Linux geprüft; für andere Systeme ist es nicht geprüft." §7 trägt dafür einen Lauf: P1 bis P11 liefen als Prozess unter Linux (Zeile „Verhalten am Protokollrand …": `uname -sr` Linux 6.8.0-139-generic), `TestE2EProtokollrandTCPProbe` deckt die plattformabhängige Stelle. Ich habe P1, P2, P4, P5, P7, P9 und die Fälle von (2) und (3) ebenso als Prozess unter Linux wiederholt; alle stimmen. Der Satz sagt nicht mehr, als die Läufe zeigen („im Container" ist gestrichen), und nennt keine Version, kein Datum, kein „bisher". README trägt den Satz nicht, ohne dass das schadet. Der Lauf von `replay` als Container in der Nacharbeit trägt den Satz nicht allein (so steht es in §7). Bestätigt.

## Quantoren im Diff (Auftrag 2)

Alle `jede/jeder/alle/keine/nie/immer/nur`-Treffer der `+`-Zeilen einzeln gelesen. Die in §7 aufgeführten Zeilen stimmen mit meinen Läufen überein; folgende tragen die Läufe nicht ganz:

| Satz | Lauf | Urteil |
|---|---|---|
| „Eine Verbindung … mit Anfrage … wird als eigene Sitzung aufgezeichnet … antwortet die Datenbank ihr mit einem COPY-Datenstrom oder sendet sie ihr eine Benachrichtigung, wird sie nicht aufgezeichnet" (§4, Mehrere Verbindungen, README) | Benachrichtigung nach der letzten Anfrage der zuhörenden Verbindung: aufgezeichnet, Exit 0 | V-3 |
| „Ohne die Option … Mit der Option wird sie, außer in den Fällen unten, eine Sitzung ohne Interaktionen" | Verbindung, die nach dem Aufbau `p` oder `CopyData` sendet: keine Sitzung, auch mit Option; die Fälle unten nennen sie nicht | V-2 |
| „ebenso eine Verbindung, deren Startnachricht das Werkzeug ablehnt (eine andere Protokollversion als 3.0, ein Parameter `_pq_.`, ein nicht lesbarer Rumpf): Sie verbraucht keine Sitzung" | Wahr; die Klammer lässt Klartext bei TLS-Pflicht für `connection` weg (Inhalt (2) der Sendung) | V-1 |
| Übrige neue Quantoren (`jede angenommene Verbindung`, `keine Sitzung`, `keine Verbindung zur Datenbank`, `handelt keine andere aus`, `kein Passwort`, `bricht er selbst ab`) | Gefahren, wahr | bestätigt |

## Plan gegen Diff (Auftrag 5)

| Prüfung | Ergebnis |
|---|---|
| Größe | `git diff --numstat cf08408..HEAD -- docs/user/benutzerhandbuch.md README.md`: 78+23 und 9+1, Summe 111, unter der Marke 200 (§4); Schätzung 70 bis 140 getroffen. Metrik nach `AGENTS.md` §3.14 genannt. |
| §8 | Zwei Liefer-Punkte (Handbuch mit Tabelle, README), eine Schicht (Dokumentation). Diff-Dateien: Handbuch, README, Plan, Review-Report; kein Code, kein Test, kein Gate, nichts unter `spec/` oder `harness/`. Stimmt. |
| §1 Abgrenzung | Kein Produkt-Code, kein Hilfetext, `benutzerhandbuch-standard.md` und `abdeckung-*.md` unberührt (`git diff --name-only` zeigt sie nicht). Eingehalten. |
| §3 | Beide Dateien als `update`; stimmt. |
| §6 R1 bis R8 | R1, R2 (P1), R3, R5 (P4), R6 (P5), R8 (P6 nicht baubar → Wortlaut ohne P6) wie entschieden umgesetzt. R4: der Satz zu `replay` steht nicht im Handbuch, wie die Alternative der Zeile es zulässt. R7: kein Satz. P3 entfällt mit dem Satz. Der Wortlaut zu `PGR-E6002` nennt die Option des Treibers erst nach P7. |
| §7 | Sätze, Proben und Quantoren folgen dem Diff (Nacharbeit zu F-1 und F-2 trägt). Abweichung: §6 und §8 nennen `BEO-REPO/verhalten-nur-unter-linux-geprueft` mit „1×", das Register zählt 3 Dateien unter `evidence/` und den Stand „3× erreicht" (V-4). |
| Metrik der Sendung (§3.13) | Nachgezählt in §8 (Architect, zweimal); keine weitere Sendung eingegangen. |

## Sendungen (Auftrag 6)

| Prüfung | Ergebnis |
|---|---|
| Geber `slice-v1-abschluss-protokollrand` | liegt in `done/`; seine Sendung an diesen Slice (fünf Inhalte) steht in §1 und DoD, mit der Kennung des Gebers; der Nehmer nimmt an (`grep -n "slice-v1-abschluss-protokollrand"` in §1 und §4). |
| Adressen dieses Slice | Der Slice schickt nichts weiter: kein Folge-Slice, kein neuer Risiko-Ausgang *eingetreten*. `slice-v1-abschluss-verbindungsende-format` nennt ihn nur als Start-Voraussetzung (Schritt 34 vor 35), nicht als Gegenstand. |
| `CancelRequest` Länge 12 | bleibt bei `slice-v1-abschluss-cancel-ohne-schluessel` (dort Ziel, DoD, §3, `startup` und `server_test.go`); nicht im Handbuch. Kein offener Verweis dieses Slice. |
| R4/P6 nicht zugesagt | bestätigt: kein Satz zu `LISTEN`/`NOTIFY` im `replay`, keiner zu „wertet nicht aus". |

## Findings

### V-1 LOW

- `quelle`: Plan §1 Inhalt (2), §6 R3 (Zählung der TLS-Pflicht-Ablehnung als „Handbuch-Bestand"), `AGENTS.md` §3.9
- `pfad`: `docs/user/benutzerhandbuch.md` · Hinweise zum Wiedergeben, Aufzählung der Zählung; `README.md` · Absatz „Am Rand des Protokolls"
- `befund`: Der Inhalt (2) der Sendung nennt „Klartext trotz TLS-Pflicht" als eine Verbindung, die nicht zählt und im `replay` keine Session verbraucht. Das Handbuch sagt es nur für `first-request` („zählt bei `first-request` weder … noch eine unverschlüsselte, die das Werkzeug mit TLS abweist (`PGR-E6003`)"), die neue Klammer für beide Zuordnungen führt Version, `_pq_.` und Rumpf, nicht Klartext. Gefahren: `replay --tls-cert` mit `connection` und mit `first-request`, erst Klartext (`PGR-E6003`), dann zwei reguläre TLS-Verbindungen: beide erhalten ihre Session, bei beiden Zuordnungen. Der Satz ist also für `connection` wahr, aber nicht geschrieben, und R3 hielt ihn zu Unrecht für Bestand (der Bestand nennt nur `first-request`).
- `verifizierbar`: ja
- `klasse`: Inhalt der Sendung nicht vollständig geliefert (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`)
- Vorschlag: „ebenso eine Verbindung, die das Werkzeug mit TLS-Pflicht abweist (`PGR-E6003`)" in die Klammer der Zählung bei beiden Zuordnungen; der Zusatz bei `first-request` entfällt oder bleibt.

### V-2 LOW

- `quelle`: `AGENTS.md` §3.11 (Quantor „sie wird eine Sitzung"), Klasse von F-1
- `pfad`: `docs/user/benutzerhandbuch.md` · Hinweise zum Aufzeichnen, Bullet „Ohne die Option … Mit der Option wird sie, außer in den Fällen unten, eine Sitzung ohne Interaktionen" und die neue Aufzählung „Sendet die Anwendung nach dem Aufbau eine Nachricht `p` … oder COPY-Daten (`CopyData`), ist das `PGR-E6001`"
- `befund`: Gefahren: `record --record-empty-sessions`; Verbindung A mit Anfrage `SELECT 1;`, danach `p`; Verbindungen B (`CopyData`), C (`CopyDone`), D (`CopyFail`) nach dem Aufbau; Verbindung E mit `SELECT 5;`. Die Datei enthält genau 1 Sitzung (E). Eine Verbindung, die nach dem Aufbau eine Nachricht `p` oder `CopyData` sendet, ist also keine Sitzung, weder mit noch ohne die Option, auch nicht mit einer abgeschlossenen Anfrage davor. Das Handbuch sagt es nicht: der neue Bullet nennt Fehlerantwort, Ende der Verbindung und Exit-Code, nicht die Aufzeichnung, und die Fälle „unten" der Option führen sie nicht. §4 sagt allgemein, eine Verbindung mit abgeschlossener Anfrage werde aufgezeichnet.
- `verifizierbar`: ja
- `klasse`: Quantor ohne die Ausnahme derselben Seite (Wiederkehr von F-1)
- Vorschlag: im neuen Bullet „die Verbindung wird nicht aufgezeichnet" ergänzen und in der Aufzählung der Option nennen.

### V-3 LOW

- `quelle`: `AGENTS.md` §3.11, Klasse von F-1
- `pfad`: `docs/user/benutzerhandbuch.md` · §4 Hinweise zum Aufzeichnen (Bullet zur Sitzung), *Mehrere Verbindungen aufzeichnen und wiedergeben* (Einleitung), *Verbindungsfehler*, Treiber-Hinweise, Tabelle `PGR-E6001`; `README.md` · „sendet sie ihr eine Benachrichtigung (`LISTEN` mit `NOTIFY`), zeichnet `record` sie nicht auf (`PGR-E6001`)"
- `befund`: Sechs Stellen sagen ohne Einschränkung, eine Verbindung, der die Datenbank eine Benachrichtigung sendet, werde nicht aufgezeichnet und sei `PGR-E6001`. Gefahren: Trifft die Benachrichtigung nach der letzten Anfrage der zuhörenden Verbindung ein, sind beide Verbindungen in der Datei und das Werkzeug endet mit Exit-Code 0. Diese Ausnahme steht nur einmal, im Bullet der Hinweise zum Aufzeichnen („kein Fehler (Exit-Code 0)"). Wer §4, das README oder die Tabelle liest, liest die Aussage ohne sie.
- `verifizierbar`: ja
- `klasse`: Quantor ohne die Ausnahme derselben Seite / Plan folgt Diff nicht
- Vorschlag: „(spätestens bei der nächsten Anfrage der zuhörenden Verbindung)" oder ein Verweis auf den Bullet an den knappen Stellen; das README darf bei „nicht aufzeichnet, außer die Benachrichtigung trifft nach der letzten Anfrage ein" bleiben.

### V-4 INFO

- `quelle`: `AGENTS.md` §3.9
- `befund`: Plan §6 (Risiko 1) und §8 nennen `BEO-REPO/verhalten-nur-unter-linux-geprueft` mit „1×". Das Register trägt 3 Dateien unter `evidence/` und den Stand „3× erreicht … Der Ausgang ist nicht zugewiesen"; die Beobachtung führt seit sechs Doku-Slices „weiter offen, kein neues Auftreten … Der Zähler bleibt bei 3". Der Plan folgt dem Register nicht in der Zahl; für die Closure ändert das den Eintrag (siehe Risiko-Ausgänge).

### V-5 INFO

- `quelle`: `AGENTS.md` §3.11
- `befund`: Die Probe P4 in §7 nennt „keine Meldung `PGR-` in beiden Läufen". Mit einer Aufzeichnung aus zwei Sessions endet der Lauf mit `PGR-W2001` (zweite Session nie zugeordnet), weil bei `first-request` ein Client, der vor der ersten Anfrage abbricht, keine Session verbraucht. Der Satz im Handbuch („`replay` meldet diesen Abbruch nicht und zählt ihn nicht als Fehler (Exit-Code 0)") bleibt wahr; er sagt nichts über Warnungen zu nicht verbrauchten Sessions. Nur der Wortlaut der Probe in §7 ist enger zu lesen (Aufzeichnung mit einer Session).

### V-6 INFO

- `quelle`: `AGENTS.md` §3.11
- `befund`: Das Handbuch nennt für `require_auth` nur `password`, `md5`, `scram-sha-256` und „bei libpq 17"; `require_auth=none` gelingt (gefahren), andere Werte (`gss`, `sspi`, `oauth`) sind nicht gefahren und stehen nicht im Handbuch. Kein Befund; die Grenze ist richtig eng gefasst.

## Risiko-Ausgänge (Kandidaten, Auftrag 7)

| Risiko (§6) | Kandidat | Begründung |
|---|---|---|
| Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`) | *weiter offen* → Register | Der Slice nennt die Grenze selbst („für andere Systeme ist es nicht geprüft"), wiederholt sie also nicht stumm; das ist kein neues Auftreten. Der Eintrag steht schon bei 3× ohne zugewiesenen Ausgang (V-4); der Lese-Schritt der Welle-Closure weist ihn zu. Die Zahl „1×" in §6 und §8 ist vor der Closure zu berichtigen. |
| Der Slice liefert seinen Teil von Handbuch und README nicht vollständig (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`) | *entfallen*, nur wenn V-1 bis V-3 behoben sind | Alle fünf Inhalte stehen, vom Verifier gegen das Binary geprüft. Bleibt V-1 offen, ist es ein Auftreten (2×) und der Ausgang *weiter offen* mit Evidenz-Datei. |

Der Lerneintrag der Closure kann die Klasse „Quantor gegen die Ausnahme derselben Seite" tragen: F-1, V-2 und V-3 sind dieselbe Klasse in drei Gestalten; Merkmal: ein neuer Satz über „jede Verbindung" oder „wird aufgezeichnet" auf einer Seite, die gleichzeitig Ausnahmen führt. (Vorschlag, kein Befund.)

## Gates und Zustand

- `make build`: grün, Image-ID `9a27f080a356` wie in §7.
- `make gates-schnell`: Exit 0 (siehe DoD-Punkt 3).
- Voller `make gates`: nicht gefahren; der Diff liegt vollständig auf der Positivliste des Stop-Hooks (siehe DoD-Punkt 3).
- Aufräumen: Container `vpg` und Netz `vnet` entfernt, keine Container und Netze übrig, Scratch-Dateien im Scratchpad, Arbeitsbaum bis auf diesen Report unverändert.
