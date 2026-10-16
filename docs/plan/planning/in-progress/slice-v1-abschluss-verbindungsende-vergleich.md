# Slice slice-v1-abschluss-verbindungsende-vergleich: Verbindungsende im Antwortvergleich (Service)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-24.a` · `LH-FA-20.a` · `SPEC-034` · `SPEC-029` · `ARC-002`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12; geschnitten am 2026-10-16 beim Anspruch (Teil 1 von zwei).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Play-Service kennt mit `--compare-responses` das Ende des Servers (`LH-FA-24.a` *Ende der Aufzeichnung*): Zeigt die Aufzeichnung an der Stelle ein Ende und sendet der Server eine `error_response` `FATAL` oder `PANIC` mit demselben SQLSTATE, ist das Ende erwartet — keine Abweichung, keine Meldung, die Session endet, auch ohne `--continue-on-error` läuft die nächste; jedes andere Ende ist die Abweichung `PGR-E5004` (Art Fehler beziehungsweise, bei aufgezeichnetem Ende und normaler Antwort, Nachrichtenart). `play` fügt der Folge `connection_closed` an und wartet nicht auf das Verbindungsende; `verglichen` zählt die Interaktion bis zur Fehlerantwort; `--allow-recorded-errors` bleibt mit Vergleich ohne Wirkung. Das Verhalten ist gegen den Fake-Server im Service belegt (Tests mit Mutationen), nicht am Binary.

**Geschnitten beim Anspruch** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap, Schritt 39a und 39b): Die Prüfung des Architect schätzte den ungeteilten Plan auf 445 bis 970, Mitte rund 710 (71 % der Marke 1000), und empfahl den Schnitt *Service | Hilfetext und E2E*. Dieser Slice (39a) behält den Service; der Hilfetext von `--allow-recorded-errors`, Abnahmeszenario 16 Satz 3 gegen PostgreSQL und der Satz in `SPEC-029` Punkt 2 sind in `slice-v1-abschluss-verbindungsende-hilfetext-binary` (39b, `next/`) aufgegangen. Einzeln lieferbar: Dieser Slice setzt nur `slice-v1-abschluss-verbindungsende-einspielen` (`done/`) voraus; 39b setzt ihn voraus, nicht umgekehrt. Die Kennung `-vergleich` bleibt bei diesem Teil, weil die Sendungen der Geber sie als Adresse nennen und der größere Teil (Service, Mutationen) hier liegt; die Teile des Nehmers stehen unten je mit ihrem Verbleib.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Ende mit `--compare-responses` (hier, im Service), der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary und Abnahmeszenario 16, Satz 3 (an 39b weitergegeben). Jener Slice (Schritt 38a, seit 2026-10-16 in `done/`) liefert die Einstufung ohne Vergleich für die einfache Anfrage und das Schließen ohne `Terminate` (Port `Einspielziel` und Einspiel-Adapter); dieser setzt ihn voraus und ändert beides nicht. Das Ende mit `--compare-responses` liefert dieser Slice im Service; der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary und Abnahmeszenario 16, Satz 3, gehen an `slice-v1-abschluss-verbindungsende-hilfetext-binary` (39b; Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13). Die Extended-Interaktion ohne Vergleich liefert `slice-v1-abschluss-verbindungsende-einspielen-extended` (Schritt 38b); sie steht vor diesem Slice in der Reihenfolge, ist aber keine Voraussetzung.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen-extended`** (dort §6 (g); Sendung des Architect, 2026-10-16, `AGENTS.md` §3.13): die Extended-Interaktion, deren letzte Gruppe mit `Flush` endet (aufgezeichnetes Ende), mit Vergleich. `play` liest dort ohne `Sync` und ohne `ReadyForQuery` bis zur Antwort auf jede Client-Nachricht der Gruppe oder bis zu einer `ErrorResponse` (die Funktion `extended`, die der Geber ohne Vergleich ändert); dieser Slice liefert den Fall mit Vergleich und die Regel, nach der sich die gelesene Folge ohne `ready_for_query` mit der Aufzeichnung vergleicht (`LH-FA-24.a` *Ende der Aufzeichnung*, vierter Punkt; Zähler `verglichen`). Der Geber liefert den Pfad ohne Vergleich, nicht diesen.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen-extended`** (dort §7, Nacharbeit nach der Verifikation V-1; Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13): die drei Kommentare im Pfad mit Vergleich, die für die letzte Gruppe mit `Flush` „bis zum ReadyForQuery“ sagen, wo `LH-FA-24.a` (*Ende der Aufzeichnung*, *Zähler*) „bis zur Antwort auf jede Client-Nachricht der Gruppe“ sagt: der Kommentar von `interaktion` in `play.go` („Mit CompareResponses hält sie die Antworten bis zum ReadyForQuery“), der Kommentar von `abschluss` in `play.go` („bis zu ihrem ReadyForQuery gelesen“) und der Kommentar von `Zaehler` in `play_vergleich.go` („bis zum ReadyForQuery oder bis zur Fehlerantwort“). Der Geber ändert sie nicht, weil ihr Pfad der mit Vergleich ist und kein Test dort den Fall mit `Flush` fährt (`AGENTS.md` §3.11); den Kommentar von `Play` hat der Geber nachgezogen.

**Übernommen aus `slice-v1-abschluss-verbindungsende-einspielen-extended`** (dort §7, Nacharbeit nach der Verifikation V-3; Sendung des Planners an den Architect-Schritt dieses Slice, 2026-10-16, `AGENTS.md` §3.13, Beschluss des Nutzers vom selben Tag „ja deine Empfehlungen“): Der Einleitungssatz von `LH-FA-24.a` *Ende der Aufzeichnung* (`spec/spezifikation.md`: „Die Folge des Servers wird bis zum `ReadyForQuery` verglichen, oder, wenn der Server eine `ErrorResponse` mit dem Schweregrad `FATAL` oder `PANIC` sendet, bis zu ihr“) nennt die Ausnahme der letzten `Flush`-Gruppe nicht, die der Aufzählungspunkt nennt. Der Architect hat den Satz am 2026-10-16 vor dem Code ergänzt, dazu den Absatz zur `ErrorResponse` in der Extended-Interaktion, der dieselbe Ausnahme brauchte; dieser Slice liefert den Satz nicht als Liefer-Punkt.

**Weitergegeben an `slice-v1-abschluss-verbindungsende-hilfetext-binary`** (Planner, 2026-10-16, `AGENTS.md` §3.13; Ursprung `slice-v1-abschluss-postgres-versionen-zusage`, dort §6, Randform „Szenarien nur gegen 17“, und Verifikation I-1): Der Satz in `SPEC-029` Punkt 2 zum erwarteten Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) und die Berichtigung des Numerus im vorhandenen Satz stehen im Commit des Integrationstests zu Szenario 16 Satz 3 — und der liegt in 39b. Dieser Slice schreibt den Satz nicht (`AGENTS.md` §3.11: vorher steht er nicht); der Geber nennt `-vergleich` als Adresse, die Adresse führt einen Hop weiter.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Einstufung ohne Vergleich (einfache Anfrage und Extended-Interaktion), das Schließen ohne `Terminate` und der Einspiel-Adapter — `slice-v1-abschluss-verbindungsende-einspielen` (Vorgänger) und `slice-v1-abschluss-verbindungsende-einspielen-extended`: andere Verzweigung des Service und andere Schicht.
- Ein Verbindungsverlust ohne Fehlerantwort bei aufgezeichnetem Ende — bleibt `PGR-E4003`, keine Abweichung (`LH-FA-24.a`); ein anderer Vorgang als das Ende des Servers.
- Aufzeichnen und Wiedergeben des Endes — `-record`, `-upstream`, `-replay`.
- Der Hilfetext von `--allow-recorded-errors` und `--compare-responses`, der Test im CLI-Adapter, Abnahmeszenario 16, Satz 3, am Binary gegen PostgreSQL, der Satz in `SPEC-029` Punkt 2 mit der Berichtigung I-1 und das akzeptierte Negativ „nur gegen 17“ — `slice-v1-abschluss-verbindungsende-hilfetext-binary` (Klasse 1): andere Schicht (CLI-Adapter, Integrationstest) und der Beleg am Binary; er nimmt die Sendungen an (dort §1 und DoD mit den Kennungen der Geber).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13), dort auch die Sätze zu `FATAL`/`PANIC` und `PGR-E5004` aus `slice-v1-abschluss-antwortvergleich-doku`. Die Handbuch-Lücke bis dahin (der Satz zu `FATAL`/`PANIC` für eine handgeschriebene Aufzeichnung mit `connection_closed` stimmt nicht mehr) ist akzeptiert, Entscheidung des Nutzers vom 2026-10-12 (dort §1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Mit Vergleich ist ein Ende erwartet, wenn die Aufzeichnung an der Position eine `error_response` `FATAL` oder `PANIC` mit `connection_closed` dahinter zeigt und der SQLSTATE (Feld `C`) beider Seiten gleich ist: keine Abweichung, keine Meldung, die Session endet, die nächste läuft ohne `--continue-on-error`. Jedes andere Ende ist `PGR-E5004`, sofort festgestellt — anderer SQLSTATE, gleicher SQLSTATE wie ein aufgezeichnetes `ERROR` ohne Ende, Aufzeichnung ohne Fehlerantwort an der Stelle (aufgezeichneter SQLSTATE leer) —; mit `--continue-on-error` endet dann die Session, die nächste läuft. Zeigt die Aufzeichnung ein Ende und antwortet der Server normal (`ERROR` mit `ReadyForQuery`), vergleicht sich die Folge nach den gewöhnlichen Regeln (Test, einfache Anfrage und Extended-Interaktion). Dazu die Extended-Interaktion, deren letzte Gruppe mit `Flush` endet (Sendung aus `slice-v1-abschluss-verbindungsende-einspielen-extended`): `play` liest nur bis zur Antwort auf jede Client-Nachricht der Gruppe oder bis zu einer `ErrorResponse` und wartet nicht auf `ReadyForQuery`; ohne Fehlerantwort `FATAL`/`PANIC` ist die Abweichung die erste verschiedene Position, sonst die Position der ersten Nachricht ohne Gegenstück (Art Nachrichtenart), `verglichen` zählt die Interaktion bis zur letzten Antwort (Test im Service, ein Fall je Seite: Fehlerantwort `FATAL` gleichen SQLSTATE erwartet, normale Antwort Abweichung; Mutation: Rückkehr vor `bisBereit` nur ohne Vergleich, der Fall liest dann ein `Naechste` mehr, als die Gruppe Antworten erwartet; dazu die drei Kommentare aus der Sendung von `slice-v1-abschluss-verbindungsende-einspielen-extended` (V-1): Sie nennen für die letzte Gruppe mit `Flush` die Antwort auf jede Client-Nachricht der Gruppe statt des `ReadyForQuery`, und genau diese zwei Fälle belegen den Satz). Dazu die Zähler: `play` fügt der Folge des Servers `connection_closed` an und liest nicht über die Fehlerantwort hinaus; `verglichen` zählt die Interaktion bis zur Fehlerantwort, `abweichend` höchstens einmal je Interaktion, `eingespielt` jede gespielte; `--allow-recorded-errors` bleibt mit Vergleich ohne Wirkung, ohne dass die Kombination ein Fehler ist (Tests im Service, je Zusage eine Mutation); die Fälle laufen in den Formen einfach, Extended Sync und Extended Flush, soweit die Form sie zulässt (§3); `play` sendet nach dem Ende keine weitere Interaktion der Session, die nächste Session läuft; ein leeres `V` gilt wie ein fehlendes (Modell und Service). Der Hilfetext, der Test im CLI-Adapter und Abnahmeszenario 16 Satz 3 liefert `slice-v1-abschluss-verbindungsende-hilfetext-binary`.
- [ ] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services/play.go`, `play_vergleich.go` | update | Ende im Vergleich: erwartet oder Abweichung, `connection_closed` anfügen, Zähler; `fatal` (heute stets Abweichung); `--allow-recorded-errors` mit Vergleich ohne Wirkung; dazu die Kommentare von `interaktion` und `abschluss` in `play.go` und von `Zaehler` in `play_vergleich.go` für die letzte Gruppe mit `Flush` (Sendung aus `slice-v1-abschluss-verbindungsende-einspielen-extended`, V-1) |
| `internal/hexagon/services` (Tests) | update | Randformen aus §6 in `play_vergleich_ende_test.go`: Fälle in Tabellen, ein Läufer (`laufVergleichEnde`, `pruefeVergleichEnde`), je Fall in den Formen einfach, Extended Sync und Extended Flush, ohne und mit `--continue-on-error` und `--allow-recorded-errors`; die Formen ohne Flush (`ohneFlush()`) fahren acht der elf Fälle von `TestVergleichEndeAbweichung`, den Fall „V ERROR bei S FATAL“ von `TestVergleichEndeGewoehnlich` und die Abweichung von `TestVergleichEndeSessionEndet`, weil aufgezeichnet kein Ende keine Gruppe mit Flush hat oder das Lesen dort mit der Antwortzahl vor der Fehlerantwort endet; `TestVergleichFatal` und `TestVergleichFatalRand` (alte Regel „stets Abweichung“) gehen in diese Tabellen ein, `TestVergleichFatalFortsetzung` in die Fälle mit `--continue-on-error` und in `TestVergleichEndeSessionEndet` (eine Interaktion D hinter B wird nicht gesendet, die nächste Session läuft), `TestVergleichFatalExtendedGruppen` bleibt mit enger gefasstem Kommentar (aufgezeichnet `ERROR` ohne Ende) |
| `test/integration` (`play_vergleich_e2e_test.go`) | nur bei Rot | `TestE2EPlayVergleichFatal` (aufgezeichnet `ERROR` 42P01 an der Stelle, Server `FATAL` 57P01: Abweichung Art Fehler, Position 1, Zähler 4/4/1) bleibt nach der Regel des Endes unverändert grün (zweiter Punkt von *Ende der Aufzeichnung*); färbt Teil 1 ihn doch rot, zieht dieser Slice ihn mit dem kleinsten Eingriff nach und meldet die Zeilen; das Nachziehen der Fälle am Binary und seiner Beschreibung liefert `slice-v1-abschluss-verbindungsende-hilfetext-binary` |
| `spec/spezifikation.md` (`LH-FA-24.a` *Ende der Aufzeichnung*, Einleitungssatz) | erledigt (Architect, 2026-10-16, vor dem Code) | Ausnahme der letzten `Flush`-Gruppe im Einleitungssatz und im Absatz zur `ErrorResponse` (Sendung aus `slice-v1-abschluss-verbindungsende-einspielen-extended`, V-3); kein Liefer-Punkt dieses Slice |
| `internal/hexagon/model/response.go`, `ende_test.go`; `internal/adapters/driven/recording` (`yaml_ende_test.go`, `sqlite_lesen_ende_test.go`) | update | Ein leeres `V` gilt wie ein fehlendes (`beendetVerbindung`, Architect zu F-6, `LH-FA-25.a`): eine Zeile Code; die Zeile „leeres V gegen S FATAL“ wandert in den Tests des Modells und beider Leser von „abgelehnt“ nach „angenommen“ |
| `docs/user/abdeckung-unit.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-einspielen` (Schritt 38a) liegt in `done/` (Voraussetzung: Schließen ohne `Terminate`, Einstufung ohne Vergleich); dieser Slice ist Schritt 39a der Reihenfolge in §5 von welle-v1-abschluss, `slice-v1-abschluss-verbindungsende-hilfetext-binary` (39b) steht dahinter. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `48325eae`, Stand des Hauptzweigs nach den Plan-Commits des Schnitts, danach nicht mehr geändert; die Commits danach tragen die Basis, den Move und die Roadmap ein und ändern keinen Pfad der Metrik); **über 800** geht der Slice zurück zur Zerlegung. **Meldeschwelle 640:** Erreicht die Summe 640, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Architect in dieser Metrik (Prüfung vor dem Code, 2026-10-16; Zahlen für den ungeteilten Plan 445 bis 970, hier der Service-Teil): **305 bis 640, Mitte 470 und damit 59 % der Marke** (Produktcode 30 bis 55, Service-Tests 280 bis 560 mit den vier nachzuziehenden Bestandstests, `dupl`-Parametrisierung 0 bis 100, Tabelle 5 bis 15; Faktor Tests zu Produktcode 7,6 bis 11,7 der Geber; die Teile addieren sich nicht auf die Enden, weil nicht alle gleichzeitig am Rand liegen).
- `in-progress` → `next`: Der Einspiel-Adapter, der Port `Einspielziel` oder der CLI-Adapter muss sich ändern — zweite Schicht, zurück zur Zerlegung (der Vorgänger schloss den Adapter ab; den CLI-Adapter trägt 39b).
- `in-progress` → `open`: Der Nutzer ändert die Regel „gleiche Stelle, gleicher SQLSTATE“ — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; „gleiche Stelle = gleicher SQLSTATE“ und „`--allow-recorded-errors` ohne Vergleich erlaubt ein aufgezeichnetes Ende“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-24.a` *Ende der Aufzeichnung*, *Meldung*, *Randformen des Vergleichs*, `LH-FA-20.a` (Tabelle der Fehlerregeln).

- *Ende bei gleichem SQLSTATE:* erwartet — keine Abweichung, keine Meldung, die Session endet, die nächste läuft auch ohne `--continue-on-error` — `LH-FA-24.a` *Ende der Aufzeichnung*.
- *Ende bei anderem SQLSTATE:* Abweichung `PGR-E5004`, Art Fehler, sofort; ohne `--continue-on-error` Abbruch, mit ihm endet die Session und die nächste läuft — `LH-FA-24.a`, Tabelle der Fehlerregeln.
- *Aufgezeichnet `ERROR` (ohne Ende), der Server sendet `FATAL` mit gleichem SQLSTATE:* Abweichung, Art Fehler — `LH-FA-24.a` *Ende der Aufzeichnung* (zweiter Punkt).
- *Ende ohne Vorbild (die Aufzeichnung hat an der Stelle keine `ErrorResponse`) und ohne `FATAL` in der Aufzeichnung:* Abweichung, aufgezeichneter SQLSTATE leer, Art Fehler — ebenda.
- *Aufgezeichnetes Ende, der Server antwortet normal:* die Folge vergleicht sich nach den gewöhnlichen Regeln, bei `ERROR` mit gleichem SQLSTATE ist die Abweichung die Position des `connection_closed` mit der Art Nachrichtenart — ebenda (dritter Punkt).
- *`FATAL` in einer früheren Interaktion als der aufgezeichneten:* Abweichung dort, die erste Abweichung gilt — `LH-FA-24.a` *Meldung* (je Interaktion die erste).
- *`--continue-on-error`:* setzt nach einer Abweichung fort; ein erwartetes Ende bricht nicht ab — Tabelle der Fehlerregeln.
- *`--allow-recorded-errors` mit Vergleich:* ohne Wirkung, ohne dass die Kombination ein Fehler ist — `LH-FA-20.a` Optionstabelle.
- *Verbindungsverlust ohne Fehlerantwort bei aufgezeichnetem Ende:* `PGR-E4003`, keine Abweichung — `LH-FA-24.a` *Ende der Aufzeichnung*.
- *Zähler:* `verglichen` zählt bis zum `ReadyForQuery`, bis zur Fehlerantwort `FATAL`/`PANIC` oder, bei der letzten Gruppe mit `Flush` eines Endes, bis zur Antwort auf jede ihrer Client-Nachrichten, `abweichend` höchstens einmal je Interaktion, `eingespielt` jede — `LH-FA-24.a` *Meldung*.
- *Extended:* Ende in der letzten Gruppe; der Vergleich der Gruppen davor bleibt, ab der Fehlerantwort gilt die Regel des Endes — `LH-FA-24.a`. **Bindung an `slice-v1-abschluss-verbindungsende-einspielen-extended` (g)** (Sendung des Architect, 2026-10-16, `AGENTS.md` §3.13): Endet die letzte Gruppe mit `Flush` (nur bei aufgezeichnetem Ende), liest `play` bis zur Antwort auf jede Client-Nachricht der Gruppe oder bis zu einer `ErrorResponse` und wartet nicht auf `ReadyForQuery`; die Funktion dafür (`extended`) ist die des Pfades mit Vergleich. Wie sich die gelesene Folge ohne `ready_for_query` mit der Aufzeichnung vergleicht, hat der Architect am 2026-10-16 in `LH-FA-24.a` *Ende der Aufzeichnung* (vierter Punkt) entschieden: die gewöhnlichen Regeln, erste verschiedene Position, sonst Position der ersten Nachricht ohne Gegenstück, Art Nachrichtenart (nicht ausschließlich die Position von `connection_closed`); `play` fügt `connection_closed` nur nach `FATAL`/`PANIC` an; `verglichen` zählt bis zur letzten Antwort. Dieser Slice hat dafür je einen Fall im Service (DoD Punkt 1, §1) und die Mutation „Rückkehr vor `bisBereit` nur ohne Vergleich“. Zur Bindung gehören die drei Kommentare, die für diese Gruppe „bis zum ReadyForQuery“ sagen (§1, DoD Punkt 1; V-1 der Verifikation von `slice-v1-abschluss-verbindungsende-einspielen-extended`): Sie werden zusammen mit dem Fall geändert, den dieser Slice liefert; bleiben sie stehen, sagt der Kommentar mehr zu als der Test prüft.
- *Stelle von `connection_closed` im Vergleich* (Prüfung des Architect, 2026-10-16, `AGENTS.md` §3.12): Auf der Seite der Aufzeichnung stellt das Format nur eine Stelle dar — letzte Nachricht, unmittelbar nach einer `error_response` `FATAL` oder `PANIC`, letzte Interaktion der Session, bei Extended die letzte Gruppe, die mit `Sync` oder `Flush` endet; Anfang, Mitte, vor dem Ende, nicht letzte Gruppe und andere Schweregrade sind eine beschädigte Aufzeichnung (`PGR-E3003`, `SPEC-002`, `SPEC-041`) und kein Fall des Vergleichs. Auf der Seite des Servers entsteht `connection_closed` nur durch das Anfügen von `play` hinter der Fehlerantwort `FATAL` oder `PANIC` und ist die letzte Nachricht der Folge. Der Vergleich sagt dazu nichts zu, was das Format nicht darstellt — `LH-FA-24.a` *Ende der Aufzeichnung*, erster Satz.
- *Reihenfolge der Gründe:* Es gibt einen Grund der Abweichung, `PGR-E5004`; die erste verschiedene Position gilt, eine Fehlerantwort `FATAL`/`PANIC` an früherer Position als die aufgezeichnete ist die Abweichung dort (Art Fehler, aufgezeichneter SQLSTATE leer, wenn dort keine `error_response` steht). Ein Verbindungsverlust ohne Fehlerantwort geht vor jedem Vergleich: `PGR-E4003`, die Interaktion wird nicht verglichen und nicht gezählt — `LH-FA-24.a` *Ende der Aufzeichnung* (Schluss), *Zähler*.
- *Server früher, gleich, später als die Aufzeichnung:* erster Unterschied der normalisierten Folge mit `connection_closed` hinter der Fehlerantwort; kein eigener Fall außer den vier Punkten — ebenda. *Schweregrad:* `FATAL` und `PANIC` gelten beidseitig als dasselbe Ende (erster Punkt, „`FATAL` oder `PANIC`“); *SQLSTATE:* zwei Fehlerantworten ohne Feld `C` sind gleich, das Ende ist dann erwartet — *Randformen des Vergleichs*, *Fehler*.
- *Einleitungssatz von `LH-FA-24.a` (Bindung an `slice-v1-abschluss-verbindungsende-einspielen-extended`, V-3; Sendung an den Architect, 2026-10-16, `AGENTS.md` §3.13):* Der Satz „Die Folge des Servers wird bis zum `ReadyForQuery` verglichen, oder … bis zu ihr“ nennt die letzte `Flush`-Gruppe nicht (Aufzählungspunkt: bis zur Antwort auf jede Client-Nachricht der Gruppe, ohne `ReadyForQuery`). Der Architect hat ihn am 2026-10-16 vor dem Code in der Spezifikation ergänzt (Einleitung und Absatz zur `ErrorResponse`, beide mit Zeiger auf den vierten Punkt); der Satz gehört zum Abgleich der Sätze einer Stelle (Einleitung, Aufzählung), den `slice-harness-regeln-spezifikation-satzabgleich` als Regel liefert. Er sagt nur zu, was die zwei Fälle der letzten `Flush`-Gruppe prüfen (DoD Punkt 1, `AGENTS.md` §3.11).
- *Zwischenstand nach diesem Slice (Schnitt, Planner, 2026-10-16):* Das Binary zeigt das Ende im Vergleich, der Hilfetext von `--allow-recorded-errors` nennt es erst mit 39b; er sagt bis dahin nichts Unwahres (er beschreibt `ERROR` ohne Vergleich), und kein Test dieses Slice fährt das Binary mit Ende und Vergleich — das Verhalten am Binary ist bis 39b nur im Service belegt (`AGENTS.md` §3.11: Der Plan sagt am Binary nicht mehr zu, als der Service belegt). Die Handbuch-Lücke (§1) beginnt mit diesem Slice.
- *Leeres `V` und `C` an anderer Nachricht* (Architect, 2026-10-16, Rückgabe zu F-6 und F-3 des Reviews): Ein leeres `V` gilt wie ein fehlendes, in `beendetVerbindung` (Modell) wie in `fehlerantwort` (Service) — `LH-FA-25.a`, Randformen des Endes in der Datei; das Modell wird angepasst (eine Zeile, Zeile „query leeres V gegen S FATAL“ in `ende_test.go` dreht um). Der SQLSTATE des Vergleichs zählt nur an einer `error_response`; ein Feld `C` an `command_complete` lässt das Format zu, der Typ-Check in `erste` ist ein Wächter mit Fall und Mutation — `LH-FA-24.a` *Fehler*.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; das Vorbild `slice-v1-abschluss-antwortvergleich` maß rund 1450 Zeilen, davon etwa 1100 Tests): vor dem Anspruch geschnitten (Schätzung des ungeteilten Plans 445 bis 970); die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Szenario 16 Satz 3 stützt sich auf `pg_terminate_backend` und `57P01`; andere Serverversionen sind nicht geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`). Das Risiko und das akzeptierte Negativ (Entscheidung des Nutzers vom 2026-10-12: nur gegen PostgreSQL 17, die Matrix fährt es nicht mit) gehören dem Test am Binary und gehen mit ihm an `slice-v1-abschluss-verbindungsende-hilfetext-binary` (dort §1 und §6, `AGENTS.md` §3.13); dieser Slice fährt keinen Server — **Ausgang:** vorgemerkt *entfallen* (kein Test gegen einen Server in diesem Slice; das Risiko steht in §6 von 39b).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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

### Belege des Implementers

**Stand:** Basis `48325eae`, Lieferung `85434400` (Code, Tests, `docs/user/abdeckung-unit.md`, Plan); Belege in diesem Abschnitt kommen danach. Keine Rückführung, keine Randform außerhalb von §6, kein Eingriff in `test/integration` (`TestE2EPlayVergleichFatal` blieb grün und unberührt), kein Port- und kein Adapter-Eingriff, Hilfetext, CLI-Test, Szenario 16 und `SPEC-029` unberührt.

**Läufe:** `make abdeckung` (7 Zeilen mehr, 5 weniger in `abdeckung-unit.md`), `make abdeckung-check` Exit 0, `make kopf-check` Exit 0, `make lint` `0 issues` (vor dem Commit), `make gates` ein Lauf, Exit 0, 12 min 48 s, kein `FAIL` im Log, `TestE2EPlayVergleichFatal` `PASS`, Nachweis `.harness/state/gates-passed.*` geschrieben.

**Wege.** Unit-Mutanten: Kopie von `go.mod`, `go.sum`, `internal`, `cmd` in ein frisches Verzeichnis je Mutant, genau eine Ersetzung (Treffer erzwungen) in `play.go` oder `play_vergleich.go`, `touch`, `docker run --rm --network=none` im Image `deps` (`docker build --target deps`), `go test -count=1 -v ./internal/hexagon/services/`; der Arbeitsbaum blieb unberührt, keine E2E-Mutanten nötig (alle Zusagen liegen im Service). „Rot“ nennt die Fälle von `TestVergleichEnde…`, die rot wurden (Formen: e = einfach, S = Extended Sync, F = Extended Flush; jeder Fall in jeder Form ohne und mit `--continue-on-error` und `--allow-recorded-errors`); `fatalFaelle` heißt vier Fälle (V FATAL, S FATAL ohne V, V PANIC, S PANIC ohne V).

**Mutationen** (Stellen: `fatal`, `fehlerantwort`, `interaktion`, `extended` in `play.go`; `erste`, `nachrichtenUnterschied` in `play_vergleich.go`):

| Zusage (Wortlaut `LH-FA-24.a`) | Mutation | roter Test |
|---|---|---|
| `play` fügt der Folge `connection_closed` an (nur dann ist das Ende erwartet) | M01: `fatal` ohne `l.v.nimm(connection_closed)` | `Erwartet`: alle sechs Fälle, e S F |
| Ende erwartet, wenn der SQLSTATE gleich ist | M02: `erste` gibt bei Art Fehler an der Position der Fehlerantwort nil | `Abweichung/anderer SQLSTATE` e S F |
| Ende erwartet, wenn die Aufzeichnung das Ende an dieser Position zeigt | M03: `erste` gibt bei Abweichung hinter der Fehlerantwort nil | `Abweichung/aufgezeichnet ERROR ohne Ende, Server` (vier Fälle, e S); `TestVergleichFatalExtendedGruppen` |
| Aufzeichnung ohne Fehlerantwort an der Stelle: Art Fehler, aufgezeichneter SQLSTATE leer; ebenso Server früher als das aufgezeichnete Ende | M05: Umwandlung nur, wenn dort eine `error_response` steht | `Abweichung/ohne Fehlerantwort der Aufzeichnung` e S; `Abweichung/Server früher` e S F |
| die erste verschiedene Position gilt; Server später als das aufgezeichnete Ende | M04: Umwandlung auch für frühere Positionen | `Abweichung/Server später` e S; `Abweichung/frühere Abweichung` e S |
| Art Fehler, nicht Nachrichtenart | M27: `art: artNachricht` bei der Umwandlung | `Abweichung/anderer SQLSTATE`, `…/Server früher`, `…/ohne Fehlerantwort der Aufzeichnung`, `…/aufgezeichnet ERROR ohne Ende, Server` (vier); `TestVergleichFatalExtendedGruppen` |
| Position der Fehlerantwort (nicht die des `connection_closed`) | M28: `position: u.position` | `Abweichung/aufgezeichnet ERROR ohne Ende, Server` (vier, e S) |
| SQLSTATE des Servers in der Meldung | M29: `sqlstateIs: ""` | `Abweichung/anderer SQLSTATE`, `…/Server früher`, `…/ohne Fehlerantwort…`, `…/aufgezeichnet ERROR ohne Ende` (vier) |
| `FATAL` und `PANIC` gelten beidseitig als dasselbe Ende (Aufzählung: aufgezeichnet PANIC, Server FATAL; Server PANIC, aufgezeichnet FATAL) | M13: der Vergleich der Fehlerantworten nimmt das Feld `V` hinzu | `Erwartet/aufgezeichnet PANIC` e S F, `Erwartet/Server V PANIC` e S F; dazu `Erwartet/Server S … ohne V` (Feld `V` fehlt), `Gewoehnlich/aufgezeichnetes Ende, Server ERROR gleicher SQLSTATE`; `TestVergleichStruktur` (zwei Fälle) |
| zwei Fehlerantworten ohne Feld `C` sind gleich, das Ende ist erwartet | M32: leerer SQLSTATE zählt als verschieden | `Erwartet/ohne Feld C` e S F; `TestVergleichStruktur/Fehler ohne Feld C auf beiden Seiten` |
| Server `PANIC` ist Ende (Aufzählung, Fall V PANIC und S PANIC ohne V) | M14: `fatal := schwere == "FATAL"` | `Erwartet/Server V PANIC`, `…/Server S PANIC ohne V` (e S F), `Erwartet/ohne Feld C`, `Abweichung/aufgezeichnet ERROR ohne Ende, Server V PANIC` und `S PANIC ohne V` (e S); dazu Bestandstests ohne Vergleich |
| Schweregrad aus `S`, wenn `V` fehlt (Fall S FATAL ohne V, S PANIC ohne V) | M15: kein Rückgriff auf `S` | `Erwartet/Server S FATAL ohne V`, `…/Server S PANIC ohne V` (e S F), `Abweichung/aufgezeichnet ERROR ohne Ende, Server S FATAL ohne V` und `S PANIC ohne V` (e S); dazu Bestandstests ohne Vergleich |
| `V` geht `S` vor (V FATAL bei S SCHWER; V ERROR bei S FATAL ist kein Ende) | M16: `S` zuerst, dann `V` | `Erwartet/Server V FATAL`, `…/Server V PANIC` (e S F), `Abweichung/aufgezeichnet ERROR ohne Ende, Server V FATAL` und `V PANIC` (e S), `Gewoehnlich/V ERROR bei S FATAL` e S |
| nichts danach: `play` liest nach der Fehlerantwort nichts mehr (verneinende Zusage; der Server hat noch zwei Antworten bereit, der Ablauf zählt die `Naechste`) | M10: `fatal` gibt beim erwarteten Ende nil zurück | `Erwartet`: sechs Fälle, e S F |
| erwartetes Ende: keine Meldung, Exit-Code 0 | M11: `fatal` gibt beim erwarteten Ende `PGR-E4003` zurück | `Erwartet`: sechs Fälle, e S F |
| `play` schließt die Verbindung ohne `Terminate` (`SchliesseOhneTerminate`) | M24: nur ohne Vergleich | `Erwartet` sechs Fälle e S F, `Abweichung` alle Fälle |
| `--continue-on-error`: die Session endet, die nächste läuft (Abweichung als Meldung weitergereicht) | M09b: `fatal` reicht immer weiter | `Abweichung`: alle Fälle ohne die Option (Abbruch erwartet, die nächste Session lief) |
| ohne `--continue-on-error` bricht der Lauf bei der Abweichung ab | M09a: `fatal` gibt immer den Fehler zurück | `Abweichung`: alle Fälle mit der Option (die nächste Session lief nicht) |
| `--allow-recorded-errors` bleibt mit Vergleich ohne Wirkung | M12: `fatal` meldet bei der Option das Ende als erwartet | `Abweichung`: alle Fälle mit `allow=true` (ohne Meldung statt `PGR-E5004`); `Erwartet` mit `allow=true` (Zähler) |
| `verglichen` zählt die Interaktion bis zur Fehlerantwort | M06: kein `Verglichen++` in `fatal` | `Erwartet` sechs Fälle, `Abweichung` alle Fälle; `TestVergleichFatalExtendedGruppen` |
| `abweichend` zählt höchstens einmal je Interaktion | M07: `Abweichend += 2` in `fatal` | `Abweichung`: alle Fälle; `TestVergleichFatalExtendedGruppen` |
| ein erwartetes Ende zählt nicht als abweichend | M08: `Abweichend++` vor der Prüfung auf nil | `Erwartet`: sechs Fälle, e S F |
| zeigt die Aufzeichnung ein Ende und antwortet der Server normal: die gewöhnlichen Regeln (einfache Anfrage und Extended) | M18: `interaktion` vergleicht bei aufgezeichnetem Ende nicht | `Gewoehnlich/aufgezeichnetes Ende, normale Antwort`, `…Server ERROR gleicher SQLSTATE`, `…Server ERROR anderer SQLSTATE` e S F |
| Verbindungsverlust ohne Fehlerantwort geht vor jedem Vergleich (`PGR-E4003`, nicht verglichen, nicht gezählt), auch bei aufgezeichnetem Ende | M22: Verbindungsverlust bei aufgezeichnetem Ende beendet die Session ohne Fehler | `Verbindungsverlust` e S F; dazu `Abweichung/anderer SQLSTATE`, `…/Server früher`, `…/Server später` (Folge der Mutation an `interaktion`) |
| letzte `Flush`-Gruppe: Rückkehr vor `bisBereit` auch mit Vergleich, `play` liest nur bis zur Antwort auf jede Client-Nachricht | M19: Rückkehr vor `bisBereit` nur ohne Vergleich | `Gewoehnlich/aufgezeichnetes Ende, …` alle drei Fälle, F (ein `Naechste` mehr im Ablauf) |
| letzte `Flush`-Gruppe: das Lesen endet mit der Fehlerantwort `ERROR` (Fall `ERROR` statt `FATAL`) | M20: `n = 0` entfällt | `Gewoehnlich/aufgezeichnetes Ende, Server ERROR gleicher SQLSTATE` und `…anderer SQLSTATE`, F; dazu Bestandstests der Geber |
| letzte `Flush`-Gruppe: `verglichen` zählt, die Abweichung ist die erste verschiedene Position, sonst die der ersten Nachricht ohne Gegenstück (Art Nachrichtenart) | M21: `interaktion` vergleicht bei letzter `Flush`-Gruppe nicht | `Gewoehnlich/aufgezeichnetes Ende, …` alle drei Fälle, F |

Alle Zusagen des Plans (DoD Punkt 1, §6) tragen eine gesehene Mutation. Die drei Kommentare der Sendung V-1 (`interaktion`, `abschluss`, `Zaehler`) sind belegt durch M19 bis M21 (die letzte `Flush`-Gruppe liest bis zur Antwort auf jede Client-Nachricht, `verglichen` zählt); der Kommentar von `interaktion` zu `AllowRecordedErrors` durch M12, der zu „liest nichts mehr“ durch M10.

**Grüne Mutanten, eingeordnet** (Weg wie oben; E1 siehe Nacharbeit, dort rot):

- E2: Schutz `ende <= len(v.soll)` entfällt — äquivalent über die Schnittstelle (F-5; der Kommentar an der Stelle nennt die Begründung): Eine Aufzeichnung, die `Validate` besteht, endet mit `ready_for_query` oder mit `connection_closed`; die Folge des Servers kann vor der Fehlerantwort keines von beiden tragen, also reicht die aufgezeichnete Folge immer bis zur Position der Fehlerantwort. Die Grenze trägt `Validate` (`PGR-E3003` beim Start); der Schutz bleibt gegen einen Panik.
- M01 färbt nur `Erwartet`: Die Fälle „aufgezeichnetes Ende an früherer oder späterer Position“ (`Server früher`, `Server später`) und „ohne Ende in der Aufzeichnung bei Server-FATAL“ (`aufgezeichnet ERROR ohne Ende`, `ohne Fehlerantwort der Aufzeichnung`) bleiben mit und ohne angefügtes `connection_closed` gleich rot, weil die Umwandlung zur Art Fehler an der Fehlerantwort das fehlende Ende deckt. Das angefügte `connection_closed` ist über die Schnittstelle nur im erwarteten Ende beobachtbar; die vom Architect genannten zwei Fälle färben M01 deshalb nicht, die Zusage trägt `Erwartet`.

**Messung** (Metrik §4, `git diff --numstat 48325eae..85434400 -- internal test docs/user`): Summe **497** (Einfügungen 355, Löschungen 142; Marke 800 zu 62 %, Meldeschwelle 640 nicht erreicht; Schätzung des Architect 305 bis 640, Mitte 470, gemessen 27 über der Mitte). Aufschlüsselung: `play.go` 39 (27/12), `play_vergleich.go` 50 (29/21), davon Produktcode ohne Kommentare 40 Zeilen und Kommentare 48; `play_vergleich_ende_test.go` 288 (neu); `play_vergleich_test.go` 108 (4/104: drei alte Tests gelöscht, einer mit neuem Kommentar); `abdeckung-unit.md` 12 (7/5). Tests 396, Produktcode 89 mit Kommentaren, 40 ohne: Faktor 4,4 (mit Kommentaren) und 9,9 (ohne), der Architect rechnete mit 7,6 bis 11,7. Der Läufer (`laufVergleichEnde`, `pruefeVergleichEnde`) trägt alle Fälle; `make lint` meldete kein `dupl`.

**Beobachtungen:** (1) `TestVergleichFatalExtendedGruppen` blieb stehen und ist nicht in die Tabelle übergegangen: Der Fall liegt in einer früheren Gruppe mit `Flush` mit eigener Aufzeichnung, die der Läufer nicht baut; sein Kommentar sagt jetzt „aufgezeichnet `ERROR` mit demselben SQLSTATE ohne Ende“, was der Test prüft (M03 färbt ihn). (2) Der Fake liefert bei leerer Folge ein `ReadyForQuery`; ein Mutant, der über die Antworten hinaus liest, hängt deshalb nicht, sondern zeigt sich als ein `Naechste` mehr im Ablauf (M19, M20) — die Tests lesen den Ablauf, nicht den Stand des Fakes. (3) `AGENTS.md` §3.10 *Aufzählung zählt je Element*: je Element der Schweregrad-Aufzählung ein Fall und ein Mutant (M14 bis M16), je Glied der Konjunktion „Ende erwartet“ ein Fall und ein Mutant (M02, M03). (4) Die Bestandstests der Geber (`TestPlayEndeErwartet` und weitere) färben bei den Mutanten M14 bis M16 und M20 mit; sie gehören den Gebern.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — der Plan nennt Schätzung und Schnitt vorab (§4). Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Ende im Vergleich mit Zählern, Hilfetext, E2E von Szenario 16 Satz 3). Schichten: zwei (Kern mit dem Service, `ARC-002`; CLI-Adapter, `ARC-005`, nur der Hilfetext). Schätzung des Planners: 600 bis 1000 Zeilen mit Tests.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-verbindungsende-einspielen-extended` (Architect, 2026-10-16, `AGENTS.md` §3.13, Ortsangabe wie Sendung zählt): zwei Fälle im Service als Teil von DoD Punkt 1 (letzte Gruppe mit `Flush`, mit Vergleich; eingetragen in §1, DoD und §6 im selben Commit wie die Berichtigung der früheren Verweisung auf eine „Zeile *Extended*“, die es in der DoD nicht gab), kein neuer Liefer-Punkt und keine Schicht: Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter); Schätzung rund 40 bis 70 Zeilen mehr (zwei Tabellenzeilen mit Fake, Metrik von §4, Faktor 10 auf wenige Zeilen Produktcode), `docs/user` konservativ wie beim Geber als dritte Schicht markiert, bis der Nutzer entscheidet, Marke und Meldeschwelle unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-postgres-versionen-zusage` (Architect, 2026-10-15, `AGENTS.md` §3.13): ein Satz in `SPEC-029` Punkt 2 im Commit des Tests von Szenario 16 Satz 3; er gehört zu DoD-Punkt 3 und fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter); Schätzung unverändert, rund zwei Zeilen mehr.

Nachgezählt beim Eintragen der Berichtigung I-1 aus `slice-v1-abschluss-postgres-versionen-zusage` (Planner, 2026-10-15, `AGENTS.md` §3.13): ein Wort im Satz von `SPEC-029` Punkt 2, den dieser Slice ohnehin im Commit seines Tests schreibt; sie gehört zu DoD-Punkt 3 und fügt weder Liefer-Punkt noch Schicht hinzu (die Spezifikation zählt nicht als Schicht, Planung ebenso nicht). Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter); Schätzung unverändert.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-einspielen-extended` (Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-einspielen`): Die Kennung steht in der Abgrenzung (§1, Ausschluss) und gibt keine Arbeit hierher; kein Code, kein Test, keine Datei kommt hinzu. Liefer-Punkte und Schichten bleiben wie in der letzten Zählung dieses Abschnitts, die Schätzung unverändert. Die Voraussetzung (Start) bleibt Schritt 38a; Schritt 38b steht davor, ohne Voraussetzung zu sein.

Nachgezählt beim Eintragen der Kommentare aus `slice-v1-abschluss-verbindungsende-einspielen-extended` (V-1; Planner, 2026-10-16, `AGENTS.md` §3.13, beim Schließen des Gebers): Die Kommentare stehen in `play.go` und `play_vergleich.go`, die dieser Slice ohnehin ändert (§3); sie gehören zu DoD Punkt 1 und fügen weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter), `docs/user` wie bisher konservativ als dritte Schicht markiert. Schätzung rund 6 Zeilen mehr (drei Kommentare in der Metrik von §4, keine Tests, die der Fall mit `Flush` nicht ohnehin trägt); Marke und Meldeschwelle unverändert.

Nachgezählt beim Eintragen des Einleitungssatzes von `LH-FA-24.a` (Planner, 2026-10-16, `AGENTS.md` §3.13, Sendung aus `slice-v1-abschluss-verbindungsende-einspielen-extended`, V-3): Der Satz steht in `spec/spezifikation.md` und ist Wortlaut des Architect vor dem Code; er ist kein Liefer-Punkt, keine Schicht (die Spezifikation zählt nicht) und kein Pfad der Metrik von §4 (`internal test docs/user`). Liefer-Punkte weiter 3, Schichten weiter zwei (Kern mit dem Service; CLI-Adapter), `docs/user` wie bisher konservativ als dritte markiert; Schätzung unverändert (Architect: 15 bis 25 Zeilen Spezifikation außerhalb der Metrik).

Geprüft vom Architect vor dem Code (2026-10-16, Metrik `AGENTS.md` §3.14, Vorprüfung vor dem Anspruch): Die Randformen aus §6 sind vollständig und in der Spezifikation entschieden; ergänzt sind in §6 die Stelle von `connection_closed`, die Reihenfolge der Gründe und die Stellung des Servers zur Aufzeichnung samt Schweregrad und leerem SQLSTATE (alle durch vorhandene Sätze entschieden, keine neue Regel). Der Stand von `51145c0` bleibt: Zeigt die Aufzeichnung ein Ende und antwortet der Server normal, gilt der gewöhnliche Vergleich (erste verschiedene Position), nicht die Position von `connection_closed`; beide fallen nur bei einfacher Anfrage mit `ERROR` und `ready_for_query` zusammen. Die Spezifikation erhielt die Ausnahme der letzten `Flush`-Gruppe in Einleitung und im Absatz zur `ErrorResponse` (Änderungshistorie). Schätzung in der Metrik von §4, Basis noch nicht eingetragen: Produktcode 35 bis 65 (`play.go`, `play_vergleich.go`, drei Kommentare, Hilfetext); Service-Tests 280 bis 560 (zehn bis sechzehn Fälle zu 12 bis 15 Zeilen, Läufer, vier nachzuziehende Bestandstests mit Löschungen und Einfügungen); `dupl`-Parametrisierung 0 bis 100; CLI-Test 15 bis 35; E2E zwei Fälle zu 46 bis 70 und `TestE2EPlayVergleichFatal` 105 bis 180; Tabellen 10 bis 30. Summe **445 bis 970, Mitte rund 710 (71 % der Marke 1000)**; nach dem Faktor (Produktcode mal 7,6 bis 11,7) 300 bis 825, der Plan rechnet mit der höheren Zahl. Das obere Ende liegt über 900: **Empfehlung des Architect: schneiden**, Linie *Service | Hilfetext und E2E*: Teil 1 (Ende im Vergleich, Zähler, Fälle der letzten `Flush`-Gruppe, drei Kommentare, Nachziehen der Bestandstests; 1 Liefer-Punkt, 1 Schicht, `docs/user` konservativ 2; Schätzung 305 bis 640, Mitte 470; Marke 800, Meldeschwelle 640) und Teil 2 (Hilfetext, E2E Szenario 16 Satz 3, Satz in `SPEC-029`; 2 Liefer-Punkte, 1 Schicht, konservativ 2; Schätzung 130 bis 245, Mitte 190; Marke 400, Meldeschwelle 300), Teil 2 nach Teil 1. Zuschnitt und Eintrag in §3.13 gehören dem Planner.

Nachgezählt beim Schnitt beim Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Teil 1 von zwei, Nehmer des Hilfetexts, der Fälle am Binary und des Satzes in `SPEC-029`: `slice-v1-abschluss-verbindungsende-hilfetext-binary`): Die Sendungen gehen hinaus, hier kommt nichts hinzu. **Liefer-Punkte: 1** (Ende im Vergleich im Service mit Zählern, den zwei Fällen der letzten `Flush`-Gruppe und den drei Kommentaren; DoD Punkt 1). **Schichten: eine** (Kern mit dem Service, `ARC-002`); `docs/user` (nur `abdeckung-unit.md`, `make abdeckung`) konservativ als zweite markiert, die Nutzerentscheidung dazu steht aus (`slice-harness-regeln-nach-v1-abschluss`); `test/integration` bleibt unberührt, solange `TestE2EPlayVergleichFatal` grün bleibt (§3). **Schätzung** in der Metrik von §4: Architect 305 bis 640, Mitte 470 (59 % der Marke 800), das obere Ende gleich der Meldeschwelle 640. Zur Linie des Schnitts: Die Prüfung des Architect schätzte den ungeteilten Plan auf 445 bis 970, Mitte rund 710 (71 % der Marke 1000), über der Grenze von 70 % und mit dem oberen Ende nahe der Marke (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); der Schnitt *Service | Hilfetext und E2E* trennt die Schicht (CLI-Adapter, Integrationstest) und den Beleg am Binary vom Verhalten, jede Hälfte ist einzeln lieferbar (Teil 1 gegen den Fake im Service, Teil 2 am Binary gegen PostgreSQL), und beide liegen mit ihrer Mitte unter 60 % der eigenen Marke. Nachgezählt im Nehmer `slice-v1-abschluss-verbindungsende-doku`: nur die Kennung im Ziel und der Verbleib der Handbuch-Lücke (dort §8, derselbe Commit).

### Nacharbeit zum Review (`docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-vergleich.md`)

**Stand:** Lieferung `a1904b57`, Basis des Reviews `18b828ab`; Architect-Commit `ccbb1073` (F-6, F-3). Weg der Mutanten wie oben, jetzt über `./internal/hexagon/services/ ./internal/hexagon/model/ ./internal/adapters/driven/recording/` (alle drei im Image `deps`, `--network=none`); unverändert grün: 3 Pakete.

**Ausgang je Befund.**

- F-1 (MEDIUM): behoben. `TestVergleichEndeSessionEndet` (Feld `danach` im Läufer: eine Interaktion D hinter B in Session 1) mit zwei Fällen, erwartetes Ende (alle drei Formen; die Aufzeichnung ist als Datei ungültig, der Service prüft die Stelle nicht, der Leser tut es, Kommentar im Test) und Abweichung mit aufgezeichnetem `ERROR` ohne Ende (zwei Formen; der Fall von `TestVergleichFatalFortsetzung`). Die Klasse „Zusage ohne Mutation“ ist damit je Weg zum Ende belegt.
- F-2 (LOW): behoben. §3 und DoD nennen die Formen je Fall; der Kommentar von `ohneFlush` nennt den Grund. Acht der elf Fälle von `TestVergleichEndeAbweichung` laufen in zwei Formen (fünf `aufgezeichnet ERROR ohne Ende`-Fälle und `ohne Fehlerantwort der Aufzeichnung` haben kein Ende, also keine Gruppe mit Flush; `Server später` und `frühere Abweichung` erreichen die Fehlerantwort in der Form Flush nicht, das Lesen endet mit der Antwortzahl), dazu `V ERROR bei S FATAL` in `Gewoehnlich`. `TestVergleichFatalFortsetzung`: §3 nennt den Verbleib.
- F-3 (LOW): behoben. Fall `Server früher, Feld C an command_complete` (`command_complete` trägt `C: 99999`, Server sendet FATAL; Art Fehler, aufgezeichneter SQLSTATE leer). E1 (Typ-Check entfällt) ist nicht äquivalent und rot; die Äquivalenz-Aussage in der Liste oben ist gestrichen, der Kommentar von `erste` bleibt in seiner Fassung („keine Fehlerantwort: leer“) und nennt die Quelle (`LH-FA-24.a` *Fehler*).
- F-4 (INFO): keine Aktion. F-5 (INFO): Kommentar an dem Schutz in `erste`, der die Äquivalenz mit `Validate` nennt (Grenze, kein Test).
- F-6 (INFO, Architect `ccbb1073`): Code im Slice. `beendetVerbindung` liest ein leeres `V` wie ein fehlendes; die Zeile „query leeres V gegen S FATAL“ wanderte von `TestValidateEndeFehler` nach `TestValidateEndeGueltig`; die gleiche Zeile in beiden Lesern (`yaml_ende_test.go`, `sqlite_lesen_ende_test.go`) wanderte von „abgelehnt“ nach „angenommen“, weil sie dieselbe Funktion prüfen. `fatalFaelle()` trägt als fünfte Ausprägung `V leer, S FATAL`; sie läuft in allen Tests, die die Liste nutzen.

**Mutationen der Nacharbeit** (Mutanten in Kopie, wie oben):

| Zusage | Mutation | roter Test |
|---|---|---|
| nach dem Ende, erwartet oder als Abweichung, endet die Session, eine weitere Interaktion wird nicht gesendet, die nächste Session läuft (F-1) | M36: `continue` statt `return nil` nach `errSessionEnde` in `session` | `TestVergleichEndeSessionEndet/erwartetes Ende` (e S F), `…/Abweichung` (e S); sonst kein Test |
| SQLSTATE des Vergleichs zählt nur an einer `error_response` (F-3) | E1: Typ-Check in `erste` entfällt | `TestVergleichEndeAbweichung/Server früher, Feld C an command_complete` e S F; sonst kein Test |
| ein leeres `V` gilt wie ein fehlendes, Modell (F-6) | M37: `beendetVerbindung` zurück auf „vorhanden, auch leer“ | `TestValidateEndeGueltig` (Zeile leeres V), `TestUnmarshalEnde/einfach, leeres V gegen S FATAL`, `TestLesenEndeAngenommen/einfach, V leer, S FATAL` |
| dasselbe beim Einspielen, Service (F-6; Konjunktionsglied `V` leer neben fehlend, V vor S, S allein) | M38: `fehlerantwort` liest `V` nur bei fehlendem Feld als leer (Rückfall auf `S` nur ohne Feld) | `TestVergleichEndeErwartet/Server V leer, S FATAL` (e S F), `TestVergleichEndeAbweichung/aufgezeichnet ERROR ohne Ende, Server V leer, S FATAL` (e S), `TestPlayEndeErwartet` und `TestPlayEndeExtendedErwartet`, Fall `V leer, S FATAL` |

**Altes Wort:** Die Suche über `internal/` (auch `ports/`) nach den Wörtern der geänderten Kommentare („vorhanden, auch leer“, „entscheidet es allein“, „bis zum ReadyForQuery“) findet nichts Veraltetes: Die Fassung in `response.go` ist umgeschrieben, die Tests der Leser und des Modells nennen „leeres V“ nur noch bei „angenommen“.

**Messung** (Metrik §4, `git diff --numstat 48325eae..a1904b57 -- internal test docs/user`): Summe **566** (Einfügungen 407, Löschungen 159; Marke 800 zu 71 %, Meldeschwelle 640 nicht erreicht); die Nacharbeit allein 81 (58/23). Aufschlüsselung: `play_vergleich_ende_test.go` 315, `play_vergleich_test.go` 109, `play_vergleich.go` 55, `play.go` 39, `abdeckung-unit.md` 22, `response.go` 8, `ende_test.go` 6, `yaml_ende_test.go` 6, `sqlite_lesen_ende_test.go` 6.

**Lauf:** `make gates` ein Lauf nach `a1904b57`, Exit 0, 12 min 46 s, kein `FAIL` im Log, Nachweis geschrieben; `make lint` `0 issues`, `make kopf-check` Exit 0.
