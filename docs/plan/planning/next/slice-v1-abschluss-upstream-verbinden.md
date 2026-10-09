# Slice slice-v1-abschluss-upstream-verbinden: Benannte Verbindung benutzen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-17.a` · `SPEC-014` · `SPEC-034` · `ARC-005`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `--upstream`, `PGWIRE_RECORDER_UPSTREAM` und der Schlüssel `upstream` lösen den Namen einer benannten Verbindung auf, die Variablen ihrer Platzhalter werden eingesetzt (`PGR-E2005`, wenn eine fehlt), und `record` verbindet zu Host und Port dieser Verbindung.

**Übernimmt:** `slice-v1-abschluss-verbindungen-platzhalter` — dessen *Benutzen der
Verbindung* nach dem vorab benannten Schnitt aus seinem §4 (*zu groß*, eingetreten am
2026-10-09 vor dem Code, Entscheidung des Nutzers; dort §1, *Abgegeben*). Der Gegenstand kam
dorthin aus `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-konfiguration` und
`slice-v1-abschluss-betrieb`. Im Einzelnen:

- die Form `host:port` und der Name einer Verbindung bei `--upstream` und
  `PGWIRE_RECORDER_UPSTREAM`, jeder gesetzte Wert nach der Zusammenführung geprüft, auch wenn
  er nicht gilt (`PGR-E2001`); die benutzte Verbindung ist die, deren Namen der
  zusammengeführte Wert nennt;
- `sslmode=require` der benutzten Verbindung bei `record` (`PGR-E2004`);
- das Einsetzen von `${VAR}` in die zerlegte URL der benutzten Verbindung, `PGR-E2005` und
  der Port nach dem Einsetzen (`PGR-E2004`), in der Reihenfolge am Ende aus `LH-FA-17.a`;
  die Konstante `PGR-E2005` liegt seit `slice-v1-abschluss-konfigurationsdatei` ohne
  Erzeuger in der Code-Tabelle, dieser Slice erzeugt sie;
- `record` verbindet zu Host und Port der benutzten Verbindung, ein IPv6-Host wieder in
  eckigen Klammern (F-528 aus dem Review von `slice-v1-abschluss-verbindungen-platzhalter`,
  §6 unten), und der Host von Option und Umgebung folgt F-521;
- der Teil des Benutzerhandbuchs dazu: in §5 *Konfigurationsdatei* die Absätze zu
  Verbindungen, `sslmode` und Platzhaltern, in §7 *Fehlercodes* die Zeilen `PGR-E2005` und
  `PGR-E2006`, für beide Slices;
- die Randformen dazu aus §6 jenes Slice (Name bei `--upstream`, `sslmode=require` bei
  `record`, Variablen der benutzten Verbindung, Einsetzen, `config show` ohne `PGR-E2005`,
  Reihenfolge am Ende, Rückgabe 8, A7 und A8 für Option und Umgebung, A9, A12; §6 unten).

**Aufsetzen.** Der Slice setzt auf dem Laden von `slice-v1-abschluss-verbindungen-platzhalter`
auf: Dort ist jede URL zerlegt, die Form der Platzhalter, `sslmode`, der Name jeder
Verbindung und der Schlüssel `upstream` sind geprüft, und ein Klartext-Passwort ist
`PGR-E2006`. Eingesetzt wird, was das Laden zerlegt hat.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Zerlegung und Grammatik der URL, Form der Platzhalter, `$$`, `${` außerhalb einer URL, der
  Name einer Verbindung, der Schlüssel `upstream` beim Laden und `PGR-E2006` —
  `slice-v1-abschluss-verbindungen-platzhalter`; er liegt vor diesem Slice, und dieser setzt
  auf ihm auf. Die Prüfung der Form `host:port` (A7) liefert jener für den Schlüssel; dieser
  Slice wendet sie auf Option und Umgebung an.
- Die Wirkung einer Verbindung bei `play` (Benutzer und Datenbank der URL mit dem Vorrang von
  `--user` und `--database`, das Passwort aus dem eingesetzten Platzhalter, TLS nach
  `sslmode`) — `slice-v1-abschluss-einspielen` (dort §1, *Übernommen aus* diesem Slice);
  `play` gibt es erst mit ihm. Hier wird das Auflösen und Einsetzen an `record` getestet.
- TLS zum Upstream bei `record` — nicht Teil des Produkts in dieser Welle
  ([welle-v1-abschluss](../welle-v1-abschluss.md) §6); `sslmode=require` ist bei `record` ein
  Fehler.
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der Slice
  ändert den CLI-Adapter und, falls die aufgelöste Adresse es verlangt, den Bootstrap; die
  Code-Tabelle im Model bleibt unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--upstream` und `PGWIRE_RECORDER_UPSTREAM` lösen einen Namen nur bei
      genauer Übereinstimmung auf, mit Vorrang vor `host:port`; die Form `host:port` folgt
      A7; jeder gesetzte Wert wird nach der Zusammenführung geprüft, die Option vor der
      Umgebungsvariable, auch wenn er nicht gilt, und ein Wert, der weder Name noch
      `host:port` ist, ist `PGR-E2001`. Benutzt ist die Verbindung, deren Namen der
      zusammengeführte Wert nennt, auch aus dem Schlüssel `upstream`; `record` verbindet zu
      ihrem Host und Port, ein IPv6-Host wieder in eckigen Klammern (Integrationstest).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `${VAR}` wird für die benutzte Verbindung einmal in die zerlegte URL
      eingesetzt, der Wert steht unverändert in seinem Teil, nur der Port wird danach
      geprüft (`PGR-E2004`). Eine nicht gesetzte oder leere Variable der benutzten Verbindung
      ist `PGR-E2005`, die Meldung nennt die Verbindung und den Namen der ersten in der
      Reihenfolge der URL, nie einen Wert; die Variablen nicht benutzter Verbindungen und bei
      `record` die in Benutzer, Passwort und Datenbank bleiben unbeachtet; `config show`
      meldet kein `PGR-E2005`. `sslmode=require` der benutzten Verbindung ist bei `record`
      `PGR-E2004`. Am Ende gilt die Reihenfolge `--upstream`, `sslmode=require`, Variablen,
      Port (Test). Beleg in §7 für Punkt 1 und 2: je Zusage Zusage · Mutation · roter Test
      (`AGENTS.md` §3.10).
- [ ] Das Benutzerhandbuch beschreibt in §5 *Konfigurationsdatei* Verbindungen, `sslmode`
      und Platzhalter und in §7 *Fehlercodes* `PGR-E2005` und `PGR-E2006` wie geliefert, den
      Stand von `slice-v1-abschluss-verbindungen-platzhalter` eingeschlossen; die
      Abdeckungstabellen sind über `make abdeckung` nachgezogen.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` (`LH-FA-17.a`) | keine Änderung durch den Implementer | Die Randformen in §6 sind entschieden (Architect, 2026-10-08 und 2026-10-09, vor dem Code); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driving/cli` | update | Form `host:port` und Auflösen des Namens bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM`, jeder gesetzte Wert geprüft (die Art `text` von `--upstream` in `cli.go` prüft heute nur den leeren Wert); benutzte Verbindung; `sslmode=require` bei `record`; Einsetzen der Variablen in die zerlegte URL des Ladens, `PGR-E2005`, Prüfung des Ports danach; Zusammensetzen von `host:port` mit geklammertem IPv6-Host |
| `internal/bootstrap` | update, falls nötig | nur, wenn die aufgelöste Adresse anders als heute an `record` übergeben werden muss |
| `internal/adapters/driving/cli` (Unit-Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a`, je Randform aus §6 ein Fall; `record` mit Verbindungsname gegen PostgreSQL |
| `docs/user/benutzerhandbuch.md` | update, falls abweichend | §5 *Konfigurationsdatei* (Verbindungen, `sslmode`, Platzhalter) und §7 *Fehlercodes* (`PGR-E2005`, `PGR-E2006`) beschreiben den Zielstand; nachgezogen wird, was der gelieferte Stand beider Slices anders sagt |
| `docs/user/abdeckung-*.md` | update | über `make abdeckung` aus den Deklarationen der neuen Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungen-platzhalter` liegt in
`done/` (Zerlegung der URL mit Platzhaltern, Schlüssel `upstream`, `PGR-E2006`). Schritt 6 der
Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (Schnitt vom 2026-10-09).
Die Randformen aus §6 entschied der Architect am 2026-10-08 und 2026-10-09 vor dem Code; vor
dem ersten Code-Commit prüft er die Liste gegen diesen Zuschnitt und gegen die Form, in der
das Laden die zerlegte URL ablegt (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder das Einsetzen verlangt einen Umbau des Zerlegers aus
  `slice-v1-abschluss-verbindungen-platzhalter`. Schnitt dann: *Auflösen und Verbinden*
  (DoD-Punkt 1) und *Einsetzen* (DoD-Punkt 2), das Handbuch beim zweiten.
- `in-progress` → `open` (blockiert — Carveout?): `record` braucht die aufgelöste Adresse in
  einer Form, die der Bootstrap ohne Änderung im Kern nicht übergeben kann; dann zuerst die
  Entscheidung.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist. Alle zogen mit dem
Schnitt vom 2026-10-09 aus §6 von `slice-v1-abschluss-verbindungen-platzhalter` hierher; dorthin
kamen sie aus §6 von `slice-v1-abschluss-konfigurationsdatei` und
`slice-v1-abschluss-konfiguration`. Die Rückgaben entschied der Architect am 2026-10-08 vor dem
Code (`f19b440`), die Nummern sind die aus jenem Plan; A7 bis A12 entschied er bei seiner
Prüfung vom 2026-10-09 vor dem ersten Code-Commit (`c297f9b`), alle in `LH-FA-17.a`. Offen ist
keine.

*Benutzen der Verbindung*

- **Name bei `--upstream`** — nur bei genauer Übereinstimmung; ein Wert, der weder Name noch
  `host:port` ist, ist ungültig (`PGR-E2001`, aus der Datei `PGR-E2004` beim Laden);
  `LH-FA-17.a`.
- **`sslmode=require` bei `record`** — `PGR-E2004` der benutzten Verbindung; `LH-FA-17.a`.
- **Variable eines Platzhalters der benutzten Verbindung nicht gesetzt** — `PGR-E2005`;
  `LH-FA-17.a`.
- **Variable eines Platzhalters gesetzt, aber leer** — nicht gesetzt, `PGR-E2005`;
  Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`. Eine leere Variable im Port ist
  `PGR-E2005`, nicht ein ungültiger Port, weil die Variablen vor dem Port geprüft werden
  (bestätigt vom Architect am 2026-10-09).
- **Platzhalter in Teilen der URL, die `record` ignoriert** (Benutzer, Passwort, Datenbank) —
  unbeachtet; `record` löst nur Platzhalter in Host und Port auf; Entscheidung des Nutzers
  vom 2026-10-08, `LH-FA-17.a`.
- **Rekursion, Einsetzen in die URL** — einmal eingesetzt, ein eingesetzter Wert wird nicht
  erneut ausgewertet; URL vor dem Einsetzen zerlegt, Wert unverändert in seinem Teil, auch mit
  `@`, `:`, `/` oder `%` (*Geheimnisse*); Entscheidung des Nutzers vom 2026-10-08,
  `LH-FA-17.a`.
- **`config show` und `PGR-E2005`** — kommt nicht vor; `config show` setzt keine Variable ein;
  `LH-FA-17.a`.
- **Reihenfolge, Teil am Ende** — nach der Zusammenführung `--upstream` (Kommandozeile, dann
  Umgebungsvariable), `sslmode=require` der benutzten Verbindung bei `record`, die Variablen
  der Platzhalter der benutzten Verbindung (`PGR-E2005`) und danach ihr Port nach dem
  Einsetzen; `LH-FA-17.a`. Die Reihenfolge innerhalb einer URL liefert
  `slice-v1-abschluss-verbindungen-platzhalter`.
- **Rückgabe 5, Teil des Kommandos** — vom Kommando hängen nur `sslmode=require` bei `record`
  und die Variablen der Platzhalter ab; dass nicht benutzte Verbindungen mitgeprüft werden,
  liefert `slice-v1-abschluss-verbindungen-platzhalter`.
- **Rückgabe 8, eingesetzter Wert macht seinen Teil ungültig** — nur der Port wird nach dem
  Einsetzen geprüft: `PGR-E2004`, im letzten Schritt direkt nach den Variablen der benutzten
  Verbindung (`PGR-E2005`). Den Host prüft der Start nicht; ein Host mit `/` scheitert beim
  Verbindungsaufbau.
- **A3, Zusammensetzen** — `record` braucht die Adresse als `host:port`; ein Host mit `:`
  (auch ein eingesetzter) wird beim Zusammensetzen wieder geklammert.
- **F-528 Klammer in der Ablage** (aus dem Review von
  `slice-v1-abschluss-verbindungen-platzhalter`, entschieden vom Architect am 2026-10-09,
  `LH-FA-17.a` *Benannte Verbindungen* und *Geheimnisse*) — die Ablage trägt nicht, ob der
  Host in Klammern stand, und braucht es nicht: Ein Host ohne Klammern enthält nach F-521
  kein `:`, auch nicht dekodiert, ein abgelegter Host mit `:` stand also in Klammern oder kommt
  aus einer Variable. *Empfehlung an den Implementer:* zusammensetzen mit Klammern genau dann,
  wenn der Host nach dem Einsetzen ein `:` enthält (dieselbe Regel wie die Standardbibliothek
  beim Zusammensetzen von Host und Port), sonst ohne. Einen eingesetzten Host prüft der Start
  nicht: `${H}` mit `a:b` ergibt `[a:b]:5432` und scheitert beim Verbinden.
  *Test:* `[::1]` und `${H}` mit `::1` ergeben `[::1]:5432`; `h` ergibt `h:5432`.
  *Mutation:* Klammern nie gesetzt, rot über `[::1]`; immer gesetzt, rot über `h`.
- **F-521, Host bei Option und Umgebung** — dieselbe Regel wie beim Schlüssel `upstream`
  (wie geschrieben, IPv6 in Klammern mit Zone hinter `%`, sonst kein Steuerzeichen, kein
  Leerraum, keines von `@`, `:`, `/`, `?`, `#`, `[`, `]`, `%`; die Zone mit denselben
  Zeichen, nicht leer, V-122 aus der Verifikation von
  `slice-v1-abschluss-verbindungen-platzhalter`; `LH-FA-17.a`), über dieselbe
  Prüfung, die `slice-v1-abschluss-verbindungen-platzhalter` liefert. *Test:*
  `--upstream 'GEHEIM@h:5'` ist `PGR-E2001`, `PGWIRE_RECORDER_UPSTREAM='GEHEIM h:5'`
  `PGR-E2001` ohne den Wert. *Mutation:* Prüfung nur am Schlüssel, rot über beide.
- **A7 Form `host:port`** bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM` — Host nicht leer
  (`:5432` ungültig), IPv6 in eckigen Klammern, Port in der Form aus A4; wie geschrieben, ohne
  Dekodierung. Ein Wert ohne `:` ist ab V-121 nur noch ein Name oder ungültig; der Vorrang des
  Namens vor `host:port` bleibt im Text stehen und ist ohne Fall (akzeptiertes Negativ: keine
  Regel zu streichen, kein Test möglich). Die Prüfung der Form liefert
  `slice-v1-abschluss-verbindungen-platzhalter` für den Schlüssel `upstream`; dieser Slice
  benutzt sie.
- **A8 Jeder gesetzte Wert von `--upstream`, Teil nach der Zusammenführung** — Option, dann
  Umgebungsvariable (`PGR-E2001`), auch wenn sie nicht gilt. Benutzt ist die Verbindung, deren
  Namen der zusammengeführte Wert nennt. Den Schlüssel `upstream` beim Laden prüft
  `slice-v1-abschluss-verbindungen-platzhalter`.
- **A9 `sslmode=require` bei `record` in der Reihenfolge** — direkt nach `--upstream`, vor
  den Variablen der Platzhalter (*Fehler*).
- **A12 Meldung zu `PGR-E2005`** — nennt die Verbindung und den Namen der ersten nicht
  gesetzten Variable in der Reihenfolge der URL, nie einen Wert.

*Akzeptiertes Negativ der Prüfung vom 2026-10-09* (keine Folgepflicht, einmalig und harmlos):

- Ein eingesetzter Wert mit Steuerzeichen bleibt ungeprüft (*Geheimnisse*: weder dekodiert
  noch geprüft); eine Umgebungsvariable trägt kein NUL, und Benutzer und Datenbank benutzt
  erst `play`.

**Risiken:**

- Die Größe: *Schätzung des Architect vom 2026-10-09* für das *Benutzen* rund 550 bis 650
  Zeilen (Code für `host:port`, Auflösen, Einsetzen und Prüfungen am Ende rund 120, Tests mit
  den Randformen oben, Integration, Handbuch und Abdeckung); das Handbuch beschreibt dabei
  auch den Stand von `slice-v1-abschluss-verbindungen-platzhalter` — **Ausgang:** offen bis
  Closure.
- Das Einsetzen hängt an der Form, in der `slice-v1-abschluss-verbindungen-platzhalter` die
  zerlegte URL ablegt; trägt sie die Platzhalter nicht je Teil und in der Reihenfolge der URL
  (A12), muss dieser Slice den Zerleger ändern (Rückführung in §4) — **Ausgang:** offen bis
  Closure.
- Das Benutzerhandbuch beschreibt Verbindungen und Platzhalter schon im Zielstand und kann vom
  gelieferten Stand abweichen (§3) — **Ausgang:** offen bis Closure.
- Zwischen der Closure von `slice-v1-abschluss-verbindungen-platzhalter` und der dieses Slice
  nimmt `--upstream` einen Verbindungsnamen an, ohne ihn aufzulösen: `record` nimmt den Namen
  dann als Adresse und scheitert spätestens beim Verbindungsaufbau; DoD-Punkt 1 schließt die
  Lücke — **Ausgang:** offen bis Closure.

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

- **Belege zur DoD (Implementer):** <…>
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `c297f9b` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (1×) — dieser Slice ist die abgegebene
  Hälfte des Schnitts von `slice-v1-abschluss-verbindungen-platzhalter`; nach der Schätzung
  (§6) liegt er unter der Grenze einer Review-Sitzung, und §4 nennt einen weiteren Schnitt
  nach DoD-Punkten nur als Rückführung. Kein Beleg vor dem Code.
- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — der Gegenstand ging durch vier Slices
  (`slice-v1-abschluss-betrieb`, `slice-v1-abschluss-konfiguration`,
  `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-verbindungen-platzhalter`);
  dieser Slice übernimmt nur das *Benutzen* und hat drei Liefer-Punkte. Wer hier mitnimmt, was
  §1 ausschließt, ändert den Plan.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (16×, verkörpert in `AGENTS.md` §3.12)
  und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (5×, verkörpert) — alle
  Randformen in §6 sind vor dem Code entschieden; eine weitere geht an den Architect zurück.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (17×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (23×, §3.11) — je Zusage eine Mutation,
  Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (19×, §3.9) — jede Korrektur zieht §1, §3 und §6 im
  selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — der Geber nennt diesen Slice in
  §1 unter *Abgegeben*, dieser ihn unter *Übernimmt*; der Nehmer
  `slice-v1-abschluss-einspielen` nennt diesen Slice in §1 unter *Übernommen aus*, alles im
  selben Commit.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor dem
Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
