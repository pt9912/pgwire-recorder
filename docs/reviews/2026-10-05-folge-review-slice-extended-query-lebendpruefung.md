# Folge-Review: slice-extended-query-lebendpruefung — 2026-10-05

**Review-Art:** Code und Spezifikation, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Nacharbeit zu F-350 bis F-357 aus dem [Review vom selben Tag](2026-10-05-review-slice-extended-query-lebendpruefung.md): `git diff c740b75..c9ba569` ohne `docs/plan/planning/open/slice-v1-abschluss-postgres-versionen.md` und `docs/plan/planning/in-progress/roadmap.md`, also die Commits `e2871a3` (Lastenheft), `4149080` (Spezifikation, Plan, Folge-Slices), `98f4140` (Code, Tests, Handbuch, Abdeckung, Plan) und `c9ba569` (Geschichtszeile der ADR). Der Planer-Commit `c1f1e43` liegt im Bereich und ändert eine Zeile in Plan §6 und `docs/plan/planning/welle-v1-abschluss.md`; er ist nur als Kontext gelesen. Die Änderungen von `4149080` an `slice-v1-abschluss-postgres-versionen` sind trotz des Pfad-Ausschlusses geprüft, weil Schwerpunkt 3 sie verlangt (`git show 4149080`).

**Skill:** `.harness/skills/reviewer.md` @ `c9ba569`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md` (ganz), `docs/plan/planning/open/slice-v1-abschluss-sessions.md` (Kopf, §1, §2, §3, §6), `docs/plan/planning/open/slice-v1-abschluss-postgres-versionen.md` (§1, §6, nur die Änderungen aus `4149080`)
- `spec/lastenheft.md` [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage); `spec/spezifikation.md` LH-FA-03.a, LH-FA-09.a, LH-FA-12.a, Historie
- [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md), ADR-Index §Konventionen
- `AGENTS.md` Hard Rules 3.3 bis 3.9; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/spec-randform-erst-im-review-entschieden`
- `docs/user/benutzerhandbuch.md` („Mit einem Datenbanktreiber arbeiten“, Meldungscodes), `docs/user/abdeckung-*.md`
- Vorheriger Report am Modul: [Review zu diesem Slice](2026-10-05-review-slice-extended-query-lebendpruefung.md) (F-350 bis F-359); höchste vergebene Nummer im Repo vor diesem Lauf F-359

**Ausgeführte Läufe im Repo:** `make docs-check` und `make abdeckung-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive c9ba569`) im Scratchpad. Unit-Tests in einem Container aus der Stufe `deps` des `Dockerfile` ohne Netz. Integrationstests aus einem Image der Stufe `integration` gegen die lokalen Images `postgres:14/16/17/18-alpine` in einem internen Docker-Netz, je nur mit `-test.run`. Die Sonden S2 und S3 sind ein eigener Testfall in der Kopie, nicht im Repo. Eigene Images, Container, Netze, Cache-Volume und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K0 | Kopie ohne Änderung: `gofmt -l`, `go vet -tags integration ./...`, `go test ./...` | grün, `gofmt` leer |
| N1 | `vtLeerraum` ab Hauptversion 0 | rot (`TestLebendpruefungServerversion`, `TestReplayLebendpruefungServerversion`) |
| N2 | `zuordnen` setzt `c.vt` nicht | rot (`TestReplayLebendpruefungServerversion`) |
| N3 | `c.vt` schon beim Handshake aus der Vorlage-Session | rot (`TestReplayLebendpruefungServerversion`, Fall „vor der Zuordnung“) |
| N4 | beim Laden `\v` für jede Session Leerraum | rot (dieselbe) |
| N5 | beim Laden `\v` für keine Session Leerraum | rot (dieselbe) |
| N6, N7 | `vtAbVersion` 16 bzw. 18 | rot (beide Serverversions-Tests) |
| N8 | Meldung nach dem Ende in `ClientMessage` mit `c.pos` (das frühere M9) | rot (`TestReplayLebendpruefungNummerNachDemEnde`) |
| N9 | dasselbe in `Query` | rot (`TestReplayLebendpruefungAufgezeichnet`) |
| N10 | `\v` nie Leerraum | rot (`TestLebendpruefungVertikalerTabulator`, `TestReplayLebendpruefungServerversion`) |
| N11, N12 | eingehende Lebendprüfung mit `vt` immer `true` bzw. immer `false` | rot (`TestReplayLebendpruefungServerversion`) |
| N13 | nicht lesbare Version zählt als ab 17 | rot (beide Serverversions-Tests) |
| I0 | `Lebendpruefung` gegen 17; `WiePostgres` gegen 14, 16, 18 | grün |
| IM1 | Produktion `vtAbVersion = 16`, `WiePostgres` gegen 16 und gegen 14 | **grün** (F-360) |
| IM2 | Produktion `vtAbVersion = 18`, `WiePostgres` gegen 17 | rot |
| IM3 | Testkonstante `vtAbVersion = 15`, gegen 17 bzw. 16 | **grün** gegen 17 (das Gate), rot gegen 16 |
| S2 | Ablauf `SELECT 1`, `"\v-- ping\v"`, `BEGIN`, `"\v"`, `SELECT 2`, `ROLLBACK` über record aufgezeichnet und über replay wiedergegeben, verglichen mit PostgreSQL | 16: gleich, die beiden `\v`-Anfragen erhalten den aufgezeichneten Syntaxfehler, der Status danach `I` bzw. `E`. 17: gleich, beide leer, Status `I` bzw. `T`. Mit Mutation N4 gegen 16: rot. F-350 ist damit Ende zu Ende behoben. |
| S3 | gegen 17 aufgezeichnet nur `SELECT 1`; im Replay `"\v"` vor der ersten Anfrage | PostgreSQL 17: leere Antwort. Replay: `PGR-E5001`, die Verbindung endet, `PGR-W2001` für die Session. So verlangt es LH-FA-09.a (*Serverversion*, vor der Zuordnung). Ergebnis zu F-361 und F-363. |

---

## Status der Findings aus dem Vor-Review

| ID | Status | Beleg |
|---|---|---|
| F-350 | **behoben** | LH-FA-09.a *Definition* und *Serverversion* binden `\v` an die Hauptversion in `server_version` der Session, ab 17. Vor der Zuordnung und ohne lesbare Version gilt `\v` nicht als Leerraum. Der Code folgt dem: beim Laden je Session, nach der Zuordnung je Verbindung (`replay.go:64`, `:236`). Jede der Mutationen N1 bis N7 und N10 bis N13 färbt einen Unit-Test rot. Sonde S2 zeigt den Fall aus F-350 gegen 16 Ende zu Ende richtig und fängt N4. Der Kommentar an `leerraum` sagt nichts mehr zu, was nicht geprüft wird. Die Grenze 17 stimmt mit S1 aus dem Vor-Review und mit I0 gegen 14, 16, 17 und 18 überein. Die ADR verlangt keine Änderung: Sie verweist für die „Zeichen des Leerraums“ auf die Spezifikation. Rest: F-360 und F-361. |
| F-351 | **behoben** | Der Kopf von `slice-v1-abschluss-sessions` nennt jetzt [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), LH-FA-09.a und [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md). §1 übernimmt das Gelieferte und nennt die Kopplung an `frei` und die Serverversion bei `connection`. DoD 1 nimmt die Lebendprüfung aus und nennt `connection` mit einer Session nur aus Lebendprüfungen. §3 und §6 (neues Risiko mit Test-Anker) ziehen nach. Rest: F-362. |
| F-352 | **begründet nicht umgesetzt, Begründung trägt** | Das `Schärft:`-Feld ist die Aufwärts-Deklaration der ADR (ADR-Index §Konventionen, `AGENTS.md` §3.4) und damit Inhalt, nicht Metadatum wie die Geschichtszeile. Es bei einer angenommenen ADR zu ergänzen, verstieße gegen §3.5. Eine Folge-ADR nur für die Spur stünde in keinem Verhältnis zu einem LOW. LH-FA-03.a und LH-FA-10.a folgen aus Sätzen der ADR, ohne sie zu erweitern. Die Lücke hat sich aber vergrößert: `4149080` ändert LH-FA-03.a ein zweites Mal unter [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) (Startfehler). Im Repo steht die Begründung noch nirgends, F-352 wird in keiner Datei außer dem Vor-Review genannt. Laut Übergabe des Vor-Reviews gehört sie in die Closure-Notiz §7; bis dahin ist das zulässig. |
| F-353 | **behoben** | `0031…md:68` trägt die Zeile `Accepted` mit `—`, in derselben Form wie [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und die folgenden angenommenen ADRs. `c9ba569` ändert nur diese Zeile (§3.5 gewahrt). |
| F-354 | **behoben** | Der `play`-Satz steht wieder am Ende des Absatzes zu `--session-assignment` (`spezifikation.md:488`). Der Absatz zu Lebendprüfungen sagt „für das Replay“ und „Das Einspielen … führt eine solche Session aus“ (`:490`–`:497`). Das deckt sich mit LH-FA-09.a („Record und Einspielen behandeln Lebendprüfungen wie jede andere Anfrage“). Widerspruchsfrei auch für `connection`: LH-FA-09.a nimmt Lebendprüfungen nur aus der Zählung der Interaktionen. Die Zuordnung nach `id` hängt nicht an dieser Zählung, und die n-te Verbindung erhält die Session wie jede ohne Interaktion. |
| F-355 | **behoben** | Der Nutzer hat die Ableitung bestätigt (Plan §1 *Startfehler*). LH-FA-03.a Schritt 2 sagt „ohne Session mit Interaktion“ und nennt den Fall, das Handbuch ebenso (`benutzerhandbuch.md:614`). Den Satz der ADR zu verwendbaren Aufzeichnungen deckt die Entscheidung des Nutzers. |
| F-356 | **behoben** | Beide Meldungen lauten „nach Interaktion N erwartet die Aufzeichnung keine weitere“. Damit ist eine spätere aufgezeichnete Lebendprüfung nicht mehr wörtlich falsch beschrieben. `ClientMessage` ist jetzt geprüft (N8 rot, früher M9 grün), `Query` weiter (N9 rot). |
| F-357 | **behoben** | Die Negative von LH-FA-09 nimmt auch die ausbleibende aufgezeichnete Lebendprüfung aus (`lastenheft.md:333`–`:337`), ohne Verweis nach unten. Draft, kein Versionssprung (`harness/conventions.md` §Versionierung). |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-360 | LOW | Der Integrationstest prüft die Versionsgrenze am Replay nur von oben. `lebendAblauf` sendet die Texte mit `\v` nur, wenn der Server sie als Leerraum liest. Gegen 14 bis 16 erreicht im Replay also keine Anfrage mit `\v` den Dienst. Eine zu niedrige Grenze im Produktcode bleibt gegen 16 und 14 grün (IM1). Die eigene Konstante des Tests prüft die Grenze nur gegen den Server, nicht gegen das Replay (IM3 rot gegen 16). Im Gate gegen 17 ist sie für jeden Wert bis 17 grün. Den Fall aus F-350, eine gegen 16 aufgezeichnete Anfrage mit `\v`, prüft kein Integrationstest; das zeigt erst Sonde S2. Zwei Texte sagen mehr zu. Die Deklaration sagt „das Replay folgt der Version in der Aufzeichnung“. `4149080` schreibt in `slice-v1-abschluss-postgres-versionen` §1, der Test „prüft die Texte mit `\v` gegen die Grenze, im Replay nach der Zuordnung“. Der Satz danach lautet „für jede Version stimmt die Antwort des Replays mit der des Servers überein“. Die geplante Matrix fängt eine falsche untere Grenze damit nicht. Die untere Grenze im Code hält nur der Unit-Test mit `"16.4"` und `"16.15"` (N6 rot). Deshalb LOW und nicht MEDIUM, obwohl das Register-Muster mehrfach aufgetreten ist. | `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) Negative; LH-FA-09.a *Serverversion*; `AGENTS.md` §3.9 (Folge-Slice) | `test/integration/lebendpruefung_e2e_test.go:153`–`:158`, `:216`–`:221`, `:229`–`:234`; `docs/user/abdeckung-e2e.md:26`; `docs/plan/planning/open/slice-v1-abschluss-postgres-versionen.md:39` | ja (IM1, IM3, S2) | Integrationstest prüft eine Versionsgrenze nur auf einer Seite |
| F-361 | LOW | Das Handbuch sagt: „Einen vertikalen Tabulator zählt die Wiedergabe wie PostgreSQL erst ab Version 17 zum Leerraum, nach der Version der aufgezeichneten Datenbank.“ Vor der ersten Anfrage einer Verbindung gilt das nicht. Eine Lebendprüfung mit `\v` ordnet dort eine Session zu und wird verglichen. Gegen eine Aufzeichnung von PostgreSQL 17 endet das mit `PGR-E5001`, wo PostgreSQL 17 leer antwortet (S3). Wenn keine Session mehr frei ist, endet es mit `PGR-E5003`. Die Spezifikation nennt die Ausnahme, die Benutzerdoku nicht. Ihr Satz „wie PostgreSQL … ab Version 17“ ist damit weiter als das Verhalten. | LH-FA-09.a *Serverversion*; Plan §1 *Handbuch*; Maintainability | `docs/user/benutzerhandbuch.md:462`–`:464`; `spec/spezifikation.md:402`–`:405` | ja (S3) | Benutzerdoku nennt eine Regel ohne ihre Ausnahme |
| F-362 | LOW | `letzteNummer` setzt eine Session mit mindestens einer erwarteten Interaktion voraus und indiziert sonst mit −1. So sagt es der Kommentar an der Funktion. Die Nacharbeit legt in LH-FA-12.a fest, dass bei `connection` die n-te Verbindung auch eine Session nur aus Lebendprüfungen erhält. Für das Replay ist das eine Session ohne Interaktion. Stellt diese Verbindung eine Anfrage, führen `Query` und `ClientMessage` über `c.pos >= len(…)` in `letzteNummer`. Der neue Absatz *Kopplung an `--session-assignment connection`* in `slice-v1-abschluss-sessions` und dessen neues Risiko nennen nur die Liste `frei`, nicht diese Annahme. Der Folge-Slice erbt sie, ohne dass sein Plan sie nennt. Heute ist der Pfad nicht erreichbar, weil `zuordnen` nur aus `frei` nimmt. | `AGENTS.md` §3.9 (betroffene Folge-Slices); `BEO-REPO/plan-folgt-korrektur-nicht`; LH-FA-12.a | `internal/hexagon/services/replay.go:167`–`:172`, `:135`–`:136`, `:189`–`:190`; `docs/plan/planning/open/slice-v1-abschluss-sessions.md:34`, `:103` | ja (Lesen) | Folge-Slice-Kopplung nennt nicht jede Invariante, die er bricht |
| F-363 | INFO | Zur Regel vor der Zuordnung. Ihre Begründung steht nirgends im Repo: Plan §1 verweist auf die Spezifikation, die Spezifikation begründet nicht, und das gehört dort auch nicht hin. In der Sache trägt sie. Die Version aus dem Handshake ist die der Session `frei[0]` beim Verbindungsaufbau. Die Session, die die Verbindung bei ihrer ersten Anfrage erhält, kann eine andere sein. Bei unterschiedlichen Versionen in einer Aufzeichnung gäbe die Handshake-Regel dann eine leere Antwort, wo der aufgezeichnete Server einen Fehler sendete. Das wäre wieder die Klasse von F-350. Die konservative Regel scheitert stattdessen laut (`PGR-E5001`), wie [ADR-0007](../plan/adr/0007-strict-replay.md) es für jede andere Abweichung tut. Ihr Preis ist S3. Er fällt nur bei Aufzeichnungen ab Version 17 an, für eine Lebendprüfung mit `\v` vor der ersten Anfrage. Kein beobachteter Treiber sendet sie. Eine Aufzeichnung mit gemischten Versionen ist mit einem einzigen `--upstream` je Aufzeichnung nur bei einem Server-Upgrade während der Aufzeichnung möglich. Der Fall, gegen den die Regel schützt, ist also ebenso selten wie ihr Preis. N3 zeigt, dass ein Test die Regel trägt. Die Begründung gehört in die Closure-Notiz, sonst entscheidet der Folge-Slice zu `connection` dieselbe Frage neu. | LH-FA-09.a *Serverversion*; [ADR-0007](../plan/adr/0007-strict-replay.md); `BEO-REPO/spec-randform-erst-im-review-entschieden` | `spec/spezifikation.md:399`–`:405`; `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md:40` | ja (N3, S3) | — (kein Fehlermuster) |

## Antwort auf die Schwerpunkte

1. **Parser für `server_version`.** `vtLeerraum` liest die führenden Ziffern. `18beta1` ergibt 18, `17.2 (Debian …)` 17, `9.6.24` 9. Fehlt der Wert, ist er leer oder beginnt er mit einem anderen Zeichen, ergibt sich `false`, ebenso bei Überlauf, weil `Atoi` dort einen Fehler liefert. Alle Varianten stehen in `TestLebendpruefungServerversion`, N1 und N13 sind rot. PostgreSQL sendet `server_version` immer, und der Record übernimmt jeden `ParameterStatus` des Handshakes (`record.go:106`–`:110`). Der Fall „fehlt“ betrifft also nur von Hand geschriebene oder ältere Aufzeichnungen ohne `server_parameters`. Für sie gilt `\v` konservativ nicht als Leerraum, wie der Nutzer entschieden hat. Ohne Befund.
2. **Regel vor der Zuordnung bei `first-request`.** Die Regel ist widerspruchsfrei zu LH-FA-09.a und LH-FA-12.a, und ein Test trägt sie (N3). Bewertung der Begründung gegen die Version aus dem Handshake: F-363. Die Folge für den Nutzer fehlt im Handbuch (F-361).
3. **Version je Verbindung nach der Zuordnung.** `zuordnen` setzt `c.vt` aus der zugeordneten Session. Vorher ist `c.vt` `false`, der Nullwert. Die Prüfung in `Query` liest `c.vt` vor `zuordnen`. Eine Lebendprüfung ordnet weiter nichts zu. `ClientMessage` erkennt keine Lebendprüfungen. N2, N11 und N12 sind rot. Ohne Befund.
4. **Überspringen aufgezeichneter Lebendprüfungen nach der Version ihrer Session.** `NewReplayService` filtert je Session mit deren Version (N4, N5 rot). Gegen 16 aufgezeichnet bleibt eine Anfrage mit `\v` eine Interaktion und liefert den aufgezeichneten Fehler samt Status `E` in der Transaktion. Gegen 17 wird sie übersprungen (S2). Die Folgen für Zuordnung, nicht verbrauchte Interaktionen und `PGR-E3004` erbt der Filter unverändert, weil sie alle an der gefilterten Liste hängen. Ohne Befund.
5. **Integrationstest mit eigener Konstante 17.** Er trägt die Grenze am Server: Gegen 16 wäre eine falsche Testkonstante rot (IM3). Er trägt die Grenze am Replay von oben: Eine zu hohe Konstante im Produktcode ist gegen 17 rot (IM2). Die untere Grenze am Replay spiegelt er nur, weil gegen 14 bis 16 kein `\v` das Replay erreicht (IM1). Sie hält allein der Unit-Test. Das Gate läuft nur gegen 17, also prüft dort auch die Testkonstante nur „17 liest `\v` als Leerraum“. Siehe F-360.
6. **LH-FA-12.a zu `connection` und `play` gegen LH-FA-09.a.** Widerspruchsfrei, siehe F-354 oben. Für LH-FA-12.a ist eine Session nur aus Lebendprüfungen „für das Replay“ eine Session ohne Interaktion. LH-FA-09.a nimmt sie nur aus der Zählung für Cursor, Zuordnung `first-request` und Verbrauch. Für `connection` und `play` folgt daraus nichts Gegenteiliges. Dass LH-FA-03.a eine Aufzeichnung mit solchen Sessions auch unter `connection` mit `PGR-E3004` ablehnt, folgt aus Schritt 2. Der Folge-Slice ist an dieser Stelle über LH-FA-03.a gebunden. Die Code-Folge für `connection` steht in F-362.
7. **Hard Rule 3.9.** Der eigene Plan folgt dem Code: Kopf, §1 (*Serverversion*, *Startfehler*, *Lastenheft*), §3 (Spezifikation, Folge-Slices, `internal/hexagon/services`, `test/integration`) und §6 Risiko 5 mit Sonde und Kandidat für den Ausgang. `slice-v1-abschluss-sessions` ist vollständig nachgezogen, bis auf F-362. In `slice-v1-abschluss-postgres-versionen` sind §1 und §6 auf die gefundene Grenze nachgezogen. Der Satz zum Integrationstest sagt dort mehr zu, als der Test prüft (F-360).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. Erkennung, Version, Filter, Zuordnung und Meldungen nach dem Ende; jede Mutation N1 bis N13 rot. Befund F-362 (Kopplung). |
| Core-Reinheit | geprüft, ohne Befund. `lebendpruefung.go` importiert nur `strconv` und `strings`. |
| `test/integration` | geprüft. Läuft gegen 14, 16, 17 und 18 grün (I0). Befund F-360. |
| `docs/user/` | geprüft. Die Abdeckungstabellen entsprechen den Deklarationen (`make abdeckung-check`). Befunde F-360 (Deklaration), F-361 (Handbuch). |
| `spec/lastenheft.md` | geprüft, ohne Befund. Kein Verweis nach unten; Draft. |
| `spec/spezifikation.md` — Strata | geprüft, ohne Befund. Keine ADR-, Slice- oder Commit-Nennung im Diff; die Historienzeile nennt nur Spec-Stellen. |
| `docs/plan/adr/` — Hard Rules 3.5, 3.8 | geprüft, ohne Befund. `c9ba569` ändert nur die Geschichtstabelle; F-352 bewusst nicht umgesetzt (siehe Status). |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Befunde F-360, F-362. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft, ohne Befund. Neue Kommentare sind Zusagen und Kopplungen im Indikativ; der Kommentar an `letzteNummer` nennt seine Annahme (F-362 betrifft den Folge-Slice, nicht den Kommentar). |
| Commit-Messages `e2871a3` bis `c9ba569` | geprüft, ohne Befund. Jede nennt `slice-extended-query-lebendpruefung`, drei nennen `LH-FA-09` im Titel, `c9ba569` nennt [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md); keine nennt eine `SPEC-*`- oder `ARC-*`-Kennung. |
| Register `BEO-REPO/*` | geprüft. F-360 ist ein weiterer Fall von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, F-362 von `BEO-REPO/plan-folgt-korrektur-nicht`. F-363 berührt `BEO-REPO/spec-randform-erst-im-review-entschieden` nicht als neuer Fall, weil die Regel vor dem Review in der Spezifikation stand. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**

- Integrationstest prüft eine Versionsgrenze nur auf einer Seite (F-360; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Benutzerdoku nennt eine Regel ohne ihre Ausnahme (F-361)
- Folge-Slice-Kopplung nennt nicht jede Invariante, die er bricht (F-362; Register `BEO-REPO/plan-folgt-korrektur-nicht`)

## Verdikt

**Merge-blockierend:** nein. F-350 ist behoben, auch Ende zu Ende gegen 16 und 17 (S2). F-351 und F-353 bis F-357 sind behoben. F-352 ist begründet offen, die Begründung gehört in die Closure-Notiz.

**Übergabe:**

- F-360 an den Implementer und den Planer von `slice-v1-abschluss-postgres-versionen`, weil dessen Matrix auf diesen Test baut.
- F-361 an den Implementer.
- F-362 an Implementer und Planer von `slice-v1-abschluss-sessions`.
- F-363 ohne erwartete Aktion außer der Begründung in der Closure-Notiz.
- Bei LOW genügt Annahme oder Begründung im Plan oder in der Closure-Notiz; der Konflikt-Pfad über den Architect ist nicht nötig.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
