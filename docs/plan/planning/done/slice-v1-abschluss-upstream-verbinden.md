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
  Reihenfolge am Ende, Rückgabe 8, A7 und A8 für Option und Umgebung, A9, A12; §6 unten),
  dazu die Entscheidungen aus der Prüfung des Architect vom 2026-10-09 vor dem Code (U1 bis
  U8, §6 unten).

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
  Das Einsetzen dieses Slice gilt für jeden Teil; `record` ruft es nur für Host und Port.
  `PGR-E2005` für eine Variable in Benutzer, Passwort oder Datenbank tritt erst bei `play` ein
  und wird dort getestet (`slice-v1-abschluss-einspielen`, dort §1, *Übernommen aus* diesem
  Slice).
- TLS zum Upstream bei `record` — nicht Teil des Produkts in dieser Welle
  ([welle-v1-abschluss](../welle-v1-abschluss.md) §6); `sslmode=require` ist bei `record` ein
  Fehler.
- Im Benutzerhandbuch die Wirkung einer Verbindung bei `play` (Passwort aus dem Platzhalter,
  `PGWIRE_RECORDER_PASSWORD`, `sslmode=require` mit TLS) — `slice-v1-abschluss-einspielen`
  (dort §1, *Übernommen aus* diesem Slice, und §3); ebenso der Abschnitt `play:` in Beispiel
  und Abschnittsliste von §5 *Konfigurationsdatei*, die hier `record:` zeigen (Verifikation
  V-125); DoD-Punkt 3 beschreibt, was geliefert ist, und `play` gibt es erst mit jenem Slice
  (Review F-533).
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der Slice
  ändert den CLI-Adapter und, falls die aufgelöste Adresse es verlangt, den Bootstrap; die
  Code-Tabelle im Model bleibt unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--upstream` und `PGWIRE_RECORDER_UPSTREAM` lösen einen Namen nur bei
      genauer Übereinstimmung auf, mit Vorrang vor `host:port`; die Form `host:port` folgt
      A7; jeder gesetzte Wert wird nach der Zusammenführung geprüft, die Option vor der
      Umgebungsvariable, auch wenn er nicht gilt, und ein Wert, der weder Name noch
      `host:port` ist, ist `PGR-E2001`. Benutzt ist die Verbindung, deren Namen der
      zusammengeführte Wert nennt, auch aus dem Schlüssel `upstream`; `record` verbindet zu
      ihrem Host und Port, ein IPv6-Host wieder in eckigen Klammern (Integrationstest).
- [x] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `${VAR}` wird für die benutzte Verbindung einmal in die zerlegte URL
      eingesetzt, der Wert steht unverändert in seinem Teil, nur der Port wird danach
      geprüft (`PGR-E2004`). Eine nicht gesetzte oder leere Variable der benutzten Verbindung
      ist `PGR-E2005`, die Meldung nennt die Verbindung und den Namen der ersten in der
      Reihenfolge der URL, nie einen Wert; die Variablen nicht benutzter Verbindungen und bei
      `record` die in Benutzer, Passwort und Datenbank bleiben unbeachtet; `config show`
      meldet kein `PGR-E2005`. `sslmode=require` der benutzten Verbindung ist bei `record`
      `PGR-E2004`. Am Ende gilt die Reihenfolge `--upstream`, `sslmode=require`, Variablen,
      Port (Test). Beleg in §7 für Punkt 1 und 2: je Zusage Zusage · Mutation · roter Test
      (`AGENTS.md` §3.10).
- [x] Das Benutzerhandbuch beschreibt in §5 *Konfigurationsdatei* Verbindungen, `sslmode`
      und Platzhalter und in §7 *Fehlercodes* `PGR-E2005` und `PGR-E2006` wie geliefert, den
      Stand von `slice-v1-abschluss-verbindungen-platzhalter` eingeschlossen; die
      Abdeckungstabellen sind über `make abdeckung` nachgezogen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` (`LH-FA-17.a`) | keine Änderung durch den Implementer | Die Randformen in §6 sind entschieden (Architect, 2026-10-08 und 2026-10-09, vor dem Code); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driving/cli` | update | Hilfe von `--upstream` nennt `host:port` oder den Namen einer Verbindung (Optionstabelle `LH-FA-17.a`); Form `host:port` und Auflösen des Namens bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM`, jeder gesetzte Wert geprüft (die Art `text` von `--upstream` in `cli.go` prüft heute nur den leeren Wert); benutzte Verbindung; `sslmode=require` bei `record`; Einsetzen der Variablen in die zerlegte URL des Ladens, `PGR-E2005`, Prüfung des Ports danach; Zusammensetzen von `host:port` mit geklammertem IPv6-Host |
| `internal/bootstrap` | keine Änderung erwartet | `RecordOptions.Upstream` trägt nach dem Auflösen die zusammengesetzte Adresse `host:port`; die Zeile beim Start und `postgres.Upstream` nennen sie wie heute (U6) |
| `internal/adapters/driving/cli` (Unit-Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a`, je Randform aus §6 ein Fall; `record` mit Verbindungsname gegen PostgreSQL |
| `docs/user/benutzerhandbuch.md` | update, falls abweichend | §5 *Konfigurationsdatei* (Verbindungen, `sslmode`, Platzhalter) und §7 *Fehlercodes* (`PGR-E2005`, `PGR-E2006`) beschreiben den gelieferten Stand beider Slices; die Wirkung bei `play` (Passwort aus dem Platzhalter, `PGWIRE_RECORDER_PASSWORD`, `sslmode=require`) steht dort nicht, sie geht an `slice-v1-abschluss-einspielen` (§1); das Beispiel zeigt einen Abschnitt `record:` mit der Verbindung `lokal`, die Abschnittsliste nennt `record:` und `replay:` (V-125) |
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
  `slice-v1-abschluss-verbindungen-platzhalter`. Schnitt dann: DoD-Punkt 1 und 2 (Code) hier,
  DoD-Punkt 3 (Handbuch und Abdeckung) als eigener Slice; der ist allein lieferbar, weil er
  geliefertes Verhalten beschreibt. Ein Schnitt zwischen DoD-Punkt 1 und 2 trägt nicht: Ohne
  das Einsetzen hätte eine benutzte Verbindung mit Platzhalter in Host oder Port bei `record`
  ein Verhalten, das `LH-FA-17.a` nicht kennt (Prüfung des Architect vom 2026-10-09).
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

*Prüfung des Architect vom 2026-10-09 vor dem Code* — gegen `LH-FA-17.a` und die Ablage in
`internal/adapters/driving/cli/verbindung.go` (`verbindung` mit `host` und `port` als `teil`,
Stücke in der Reihenfolge der URL, Platzhalter als `stueck.variable`, Port ohne Angabe
`5432`, Port mit Platzhalter beim Laden auf Ziffern geprüft; `datei.namen` die gültigen Namen;
`hostPortForm` für den Schlüssel). Die Ablage trägt alles, was dieser Slice braucht: Das
Einsetzen ist ein Durchlauf über die Stücke, A12 folgt aus ihrer Reihenfolge; das zweite
Risiko unten ist damit vor dem Code entfallen. Neu entschieden heißt: seit dem Commit dieser
Prüfung in `LH-FA-17.a` (*Benannte Verbindungen*, Absatz nach *Die Wirkung einer URL*).

Bestätigt, schon entschieden:

- **Variable mit Steuerzeichen** — bleibt ungeprüft (akzeptiertes Negativ oben); in Host oder
  Port: ein Port mit Steuerzeichen hat nicht die Form eines Ports (`PGR-E2004`, Rückgabe 8),
  ein Host scheitert beim Verbinden (`PGR-E4002`).
- **Eingesetzter Host mit `:` ohne IPv6-Form** — geklammert, scheitert beim Verbinden (F-528).
- **Eingesetzter Port mit Überlauf** — nicht die Form eines Ports, `PGR-E2004` (Rückgabe 8);
  `5${P}` mit `P=5432` ergibt `55432`, gültig.
- **Mehrere fehlende Variablen** — der Start endet beim ersten Fehler; `PGR-E2005` nennt die
  erste in der Reihenfolge der URL (A12), bei `record` nur unter Host und Port.
- **Benutzer, Datenbank und Passwort bei `record`** — ignoriert, auch ihre Variablen
  (*Wirkung einer URL*, Entscheidung des Nutzers vom 2026-10-08).

Neu entschieden:

- **U1 Namen nur aus der gewählten Datei** — ohne Datei ist `--upstream staging` weder Name
  noch `host:port`, `PGR-E2001`. *Test:* ohne Datei `--upstream staging` `PGR-E2001`; mit
  Datei, die `staging` führt, gültig. *Mutation:* Name gegen eine leere Menge statt gegen
  `datei.namen` geprüft, rot über den zweiten Fall.
- **U2 Variablen einmal gelesen** — nach der Zusammenführung, eine spätere Änderung wirkt
  nicht; `record` verbindet jede Session zu derselben Adresse. *Test:* Unit, die Adresse in
  `RecordOptions.Upstream` ist nach dem Lesen fest (`os.Setenv` danach ändert sie nicht).
  *Mutation:* Einsetzen beim Verbinden statt beim Lesen — im Zuschnitt nicht baubar, ohne
  `postgres.Upstream` zu ändern (§1, Schicht-Abgrenzung); der Test sichert die Form der
  Übergabe.
- **U3 Stelle der Meldungen nach dem Einsetzen** — `sslmode=require` bei `record`,
  `PGR-E2005` und der Port nach dem Einsetzen nennen `connections.<Name>`, nie einen Wert;
  `PGR-E2005` dazu den Namen der Variable (A12). *Test:* je Fall die Meldung genau, mit
  `GEHEIM` als Wert einer Variable im Port (`PGR-E2004` ohne `GEHEIM`). *Mutation:* Wert in
  die Meldung zum Port, rot über `GEHEIM`.
- **U4 Port beim Zusammensetzen** — so, wie er geschrieben oder eingesetzt ist, auch
  `05432`; die Wählfunktion der Standardbibliothek liest führende Nullen. *Test:* `h:05432`
  ergibt die Adresse `h:05432`; Integration: `record` mit einer Verbindung, deren Port aus
  einer Variable kommt, zeichnet auf. *Mutation:* Port aus dem Teil vor dem Einsetzen
  genommen, rot über die Variable.
- **U5 Klammern** — ein Host mit `:` steht in eckigen Klammern, auch ein eingesetzter, der
  keine IPv6-Adresse ist (F-528). Test und Mutation wie bei F-528 oben; dazu `${H}` mit `a:b`
  ergibt `[a:b]:5432`.
- **U6 Adresse in Log und Meldung** — Host und Port sind kein Geheimnis: Die Zeile beim Start
  (`upstream=`) und die Meldungen zum Upstream (`PGR-E4002`) nennen die zusammengesetzte
  Adresse, nicht den Namen; nie Benutzer, Passwort oder Datenbank. *Test:* Integration,
  `record` mit `--upstream staging` (URL mit `u:${PW}@`) schreibt `upstream=<host:port>` in
  die Zeile beim Start, ohne den Wert von `PW` und ohne `u`. *Mutation:* `upstream=` mit dem
  geschriebenen Wert der Option, rot (`staging` statt der Adresse).
- **U7 Hilfe von `--upstream`** — nennt `host:port` oder den Namen einer Verbindung, wie die
  Optionstabelle. *Test:* `TestHilfe…` prüft den Text. *Mutation:* alter Text, rot.
- **U8 Einsetzen bei `play`** — das Einsetzen gilt für jeden Teil; `PGR-E2005` für eine
  Variable in Benutzer, Passwort oder Datenbank tritt erst bei `play` ein und wird in
  `slice-v1-abschluss-einspielen` getestet (dort §1 eingetragen, §3.13). Hier: Unit-Test des
  Einsetzens über alle Teile (Reihenfolge der URL, erste fehlende Variable).

*Akzeptiertes Negativ der Prüfung vom 2026-10-09* (keine Folgepflicht, einmalig und harmlos):

- Ein eingesetzter Wert mit Steuerzeichen bleibt ungeprüft (*Geheimnisse*: weder dekodiert
  noch geprüft); eine Umgebungsvariable trägt kein NUL, und Benutzer und Datenbank benutzt
  erst `play`.

**Risiken:**

- Die Größe: *Schätzung des Architect vom 2026-10-09* für das *Benutzen* rund 550 bis 650
  Zeilen; *neu geschätzt vor dem Code am 2026-10-09* gegen die Ablage: Code 150 bis 200
  (Prüfung von Option und Umgebung nach der Zusammenführung 40, Auflösen und Prüfungen am Ende
  60, Einsetzen 40, Zusammensetzen 15, Hilfe), Unit-Tests 350 bis 450, Integration 80 bis
  120, Handbuch 60 bis 80, Abdeckung 20 — **660 bis 870 Zeilen**. Der Geber lag mit 1066
  Zeilen um die Hälfte über seiner Schätzung; mit demselben Faktor rund 1000 bis 1300. Er war
  mit rund 70 Zusagen in einer Sitzung prüfbar; dieser trägt rund 30. Kein Schnitt vor dem
  Code; der Schnitt in §4 (Handbuch abtrennen) bleibt die Rückführung — **Ausgang:**
  entfallen: Der Diff `9a04b12..80cb896` hat 712 hinzugefügte Zeilen, innerhalb der Schätzung,
  und war in einer Review-Sitzung prüfbar (Review F-536); keine Rückführung trat ein.
- Das Einsetzen hängt an der Form, in der `slice-v1-abschluss-verbindungen-platzhalter` die
  zerlegte URL ablegt; trägt sie die Platzhalter nicht je Teil und in der Reihenfolge der URL
  (A12), muss dieser Slice den Zerleger ändern (Rückführung in §4) — **Ausgang:** entfallen:
  Die Ablage trägt je Teil die Stücke in der Reihenfolge der URL (Prüfung des Architect vom
  2026-10-09 vor dem Code, oben).
- Das Benutzerhandbuch beschreibt Verbindungen und Platzhalter schon im Zielstand und kann vom
  gelieferten Stand abweichen (§3) — **Ausgang:** eingetreten für §5 *Konfigurationsdatei*
  (Verifikation V-125): Das Beispiel mit dem Abschnitt `play:` und die Liste
  „(`record:`, `replay:`, `play:`)“ lehnte das gelieferte Binary mit `PGR-E2004` ab, und
  `--upstream staging` aus dem Satz danach ist bei `record` wegen `sslmode=require`
  ebenfalls `PGR-E2004`. Behoben in diesem Slice (Nacharbeit zur Verifikation, §7):
  Beispiel mit `record:` und `upstream: lokal`, Liste ohne `play:`, der Satz nennt `lokal`.
  Die Rückkehr von `play:` in Beispiel und Liste übernimmt `slice-v1-abschluss-einspielen`
  (dort §1 und §3, mit der Kennung dieses Slice).
- Zwischen der Closure von `slice-v1-abschluss-verbindungen-platzhalter` und der dieses Slice
  nimmt `--upstream` einen Verbindungsnamen an, ohne ihn aufzulösen: `record` nimmt den Namen
  dann als Adresse und scheitert spätestens beim Verbindungsaufbau; DoD-Punkt 1 schließt die
  Lücke — **Ausgang:** entfallen: DoD-Punkt 1 ist geliefert, `--upstream staging`,
  `PGWIRE_RECORDER_UPSTREAM` und der Schlüssel `upstream` lösen am Binary auf (Verifikation,
  Abschnitt 1, Punkt 1).

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

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten. DoD-Punkt 1, 2 und
  `make gates` bestätigt die Verifikation am Stand `b353056`
  (`docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-upstream-verbinden.md`,
  Abschnitt 1; Binary gegen PostgreSQL, 19 eigene Unit-Mutanten und ein Integrations-Mutant,
  alle rot), Punkt 3 mit der Einschränkung V-125. Die Nacharbeit zur Verifikation `d64f8e9`
  ist am Diff geprüft: Beispiel und Abschnittsliste in §5 *Konfigurationsdatei* zeigen
  `record:` statt `play:`, der Satz danach nennt `lokal`; §6 und §7 ziehen nach, die Sendung
  steht in `slice-v1-abschluss-einspielen` §1 und §3. Der Planner hat das Beispiel am Stand
  `d64f8e9` nachgefahren (`make build`, Beispiel per `awk` wörtlich als Datei, `docker run
  --network none`): `config show --config` Exit 0; `record --config` mit beschreibbarem
  `./recordings/` startet (`record gestartet listen=127.0.0.1:15432 upstream=localhost:5432`)
  und endet nach `docker stop` mit Exit 0; mit `--upstream staging` `PGR-E2004` wegen
  `sslmode=require`, wie der Absatz zu `record` sagt. Ohne das Verzeichnis `./recordings/`
  endet `record` mit `PGR-E3001`; das ist die Zusage zu `--output`, nicht zum Beispiel. Damit
  ist DoD-Punkt 3 ohne Einschränkung erfüllt. V-126 ist in §7 eingetragen (U2 in zwei Zeilen,
  erste Hälfte mit Mutant A1).
- **Was hat funktioniert:** Die Prüfung des Architect vor dem Code gegen die Ablage des Gebers
  (`c297f9b`) trug: Das Einsetzen war ein Durchlauf über die Stücke, kein Umbau des Zerlegers,
  keine Randform ging zurück, und weder Review noch Verifikation fanden eine Randform außerhalb
  von §6. Die Schätzung vor dem Code (660 bis 870 Zeilen) traf den Diff (712); der Slice blieb
  in einer Review-Sitzung prüfbar (F-536). Der Implementer belegte 38 Mutanten, 37 rot, den
  grünen fing ein eigener Testfall (`80cb896`); die Verifikation fuhr 20 eigene nach, alle rot.
  Die Sendungen an `slice-v1-abschluss-einspielen` (F-533, V-125) standen im selben Commit im
  Nehmer (§3.13).
- **Was ging anders als geplant:**
  1. Eine Zusage hielt nur für einen Teil ihrer Fälle: Die Auswahl der Verbindung nach dem
     genauen Namen hatte keinen roten Mutanten (F-534, S grün); §7 belegte die Zeile mit der
     Prüfung über `hatName`, die die Auswahl nicht hält. Behoben in `b353056`.
  2. Das Handbuch stand im Zielstand und wich vom gelieferten Stand ab, wie das dritte Risiko
     in §6 erwartete; gefunden hat es nicht das Lesen (Implementer, Review: „Die Zusagen zu
     `record` stimmen mit dem Code überein“), sondern der Lauf des Beispiels als Datei in der
     Verifikation (V-125). §3 sagte „nachgezogen wird, was der gelieferte Stand anders sagt“,
     nannte aber keinen Weg, die Abweichung zu finden.
  3. Zwei Kommentare sagten mehr oder weniger, als der Code tut (F-530, F-531), und §7
     notierte U2 schwächer, als der Beleg ist (V-126).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-09-review-slice-v1-abschluss-upstream-verbinden.md`: „0 HIGH · 1
    MEDIUM · 2 LOW · 4 INFO. F-530: Der Kommentar an `SetzeEin` sagt `""` für einen nicht
    geschriebenen Port; es ist `5432`. F-531: Die Kommentare an `zuletzt` nennen eine Prüfung
    und verschweigen, dass `zuletzt` die Adresse schreibt. F-532: Die Reihenfolge
    Pflichtoptionen, Kombinationen, `--upstream` ist erfüllt, weil `record` keine Kombination
    prüft. F-533: Das Handbuch nennt `play` im Zielstand; ob das zu DoD-Punkt 3 passt,
    entscheidet der Verifier. F-534: Die Auswahl der Verbindung nach dem genauen Namen hält
    keine Mutation (Mutant S grün, Sonde unterscheidet). F-535: `TestUpstreamEinmalGelesen`
    trägt nur in seiner ersten Hälfte. F-536: Der Diff war in einer Sitzung prüfbar.
    Wiederkehrende Klassen: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-534),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-530).“ Verifikation
    `docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-upstream-verbinden.md`: „0 HIGH ·
    0 MEDIUM · 1 LOW · 1 INFO. V-125: Das Beispiel in §5 *Konfigurationsdatei* wird vom
    gelieferten Binary mit `PGR-E2004` (`play: unbekannter Schlüssel`) abgelehnt; DoD-Punkt 3
    hält nur mit dieser Einschränkung. V-126: §7 notiert für U2 keine Mutation, obwohl A1 die
    erste Hälfte des Tests hält.“ Ausgänge: F-530, F-531, F-533 und F-534 in `b353056`;
    F-532 an `slice-v1-abschluss-tls-client` (Folge-Slices unten); F-535 und F-536 von der
    Verifikation eingeordnet (Grenze in §7, erstes Risiko); V-125 und V-126 in `d64f8e9`, die
    Rückkehr von `play:` an `slice-v1-abschluss-einspielen`.
- **Steering-Loop-Eintrag:** Geschärfte Regel, im Plan des Nehmers angewandt: Ein DoD-Punkt
  „Handbuch wie geliefert“ über einen Abschnitt mit Beispiel ist erst belegt, wenn das Beispiel
  wörtlich als Datei gegen das gebaute Binary läuft; Lesen gegen den Code fand die Abweichung
  nicht, der Lauf fand sie sofort (V-125). `slice-v1-abschluss-einspielen` trägt das als
  Kriterium seiner Sendung in §1 („das Beispiel als Datei startet mit `config show` ohne
  Meldung“, seit `d64f8e9`). Kein Feld `liegt in`: Repo-weit ist mit diesem Slice nichts
  verkörpert; gezählt ist der Fund unter `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
  (verkörpert in `AGENTS.md` §3.11, dessen Liste Kommentar, Hilfetext,
  Abdeckungs-Deklaration und Plan-Zeile nennt, das Handbuch nicht). Ob §3.11 oder
  `implement-slice` das Handbuch ausdrücklich aufnimmt, entscheidet der Nutzer.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist nicht wieder
  aufgetreten: Kein Code-Commit änderte den Plan, und jede Nacharbeit zog den Plan im selben
  Commit nach (`b353056` §1, §3 und §7; `d64f8e9` dazu §6); die Regel bleibt. §3.10 (seit
  welle-extended-query) ist wieder aufgetreten (F-534); die Regel bleibt, der Sensor ist mit
  `slice-harness-mutation` geplant. §3.11 (seit welle-extended-query) ist wieder aufgetreten
  (F-530, F-531, V-125); die Regel bleibt. §3.12 (seit slice-harness-randformen-vor-code) ist
  in beiden Hälften nicht wieder aufgetreten: keine Rückgabe, kein Code-Commit trägt eine
  Randform in §6 ein, und Review und Verifikation fanden keine Randform außerhalb von §6. §3.13
  (seit slice-lint-bestand-kern-driven) ist nicht wieder aufgetreten: Der Nehmer trug F-533 und
  V-125 im selben Commit, die Verifikation prüfte die Annahme. Die Regel *Beleg im Plan, nicht
  im Bericht* (`.claude/agents/implementer.md`, seit slice-harness-blackbox-kern) ist nicht
  wieder aufgetreten: §7 trägt jeden Gate-Lauf, auch nach beiden Nacharbeiten.
  `implement-slice` Schritt 19 (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`): nicht
  wieder aufgetreten, jeder Mutant in einer frischen Kopie.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `d64f8e9` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 18× → 19× (F-534), Stand
    verkörpert, bleibt.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 24× → 25× (F-530, F-531,
    V-125 mit F-533), Stand verkörpert, bleibt.
  - `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`: **kein Beleg**, bleibt 1× (offen). Dieser
    Slice ist die abgegebene Hälfte des Schnitts; seine Schätzung lag unter der Grenze, und die
    vorab benannte Rückführung trat nicht ein.
  - `BEO-REPO/slice-waechst-durch-uebernahmen`: **kein Beleg**, bleibt 2× (offen). Nach der
    Anlage nahm der Slice nichts auf; er gab das Handbuch zu `play` an den Nehmer ab.
  - `BEO-REPO/gruener-mutant-faelschlich-aequivalent`: **kein Beleg**, bleibt 1× (offen). Der
    grüne Mutant des Implementers war nicht als äquivalent eingestuft, sondern bekam einen Test
    (`80cb896`); S war nicht eingeordnet, sondern nicht gefahren, das zählt unter
    `negativtests-fehlen-bei-neuem-vertrag`.
  - Ohne Beleg, verkörpert: `BEO-REPO/plan-folgt-korrektur-nicht` (bleibt 20×),
    `BEO-REPO/spec-randform-erst-im-review-entschieden` (bleibt 17×),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 5×),
    `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 3×),
    `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (bleibt 4×); ohne Beleg, offen:
    `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (bleibt 1×).
  - Einmalig und nicht eingetragen: F-532 (latente Reihenfolge, an den Folge-Slice), F-535
    und F-536 (Hinweise), V-126 (Beleg schwächer notiert, als er ist). Mit diesem Slice
    erreicht kein Eintrag die Schwelle 3× neu; über der Schwelle stehen nur Einträge mit
    Ausgang (`commit-nennt-struktur-kennung` 3× und `roter-lauf-haengt-bis-zum-zeitlimit` 3×
    geplant; `folge-slice-adresse-nimmt-nicht-an` 3×,
    `implementer-bericht-erreicht-pruefer-nicht` 4×,
    `randform-im-code-entschieden-dann-zurueckgegeben` 5×,
    `spec-randform-erst-im-review-entschieden` 17×, `negativtests-fehlen-bei-neuem-vertrag`
    19×, `plan-folgt-korrektur-nicht` 20×, `zusage-im-kommentar-weiter-als-pruefung` 25×
    verkörpert; `white-box-liste-vor-code-nur-namenssuche` 3× gestrichen).
- **Folge-Slices:** `slice-v1-abschluss-einspielen` (`next/`; aus §1 *Ausdrücklich NICHT*: die
  Wirkung einer Verbindung bei `play`, `PGR-E2005` in Benutzer, Passwort und Datenbank nach
  U8, im Handbuch die Wirkung bei `play` nach F-533 und der Abschnitt `play:` in Beispiel und
  Liste nach V-125; dazu das dritte Risiko aus §6). `slice-v1-abschluss-tls-client` (`next/`;
  F-532: die Kombinationen von `record` vor jedem Fehler von `--upstream`, mit Test, zu seinem
  DoD-Punkt 2; mit dieser Closure in seinem §1 eingetragen). Der nächste Schritt nach §5 von
  [welle-v1-abschluss](../welle-v1-abschluss.md) ist `slice-v1-abschluss-schreiben`.
- **Risiken aus §6:** vier. *Größe*: **entfallen**, 712 Zeilen in der Schätzung, eine
  Review-Sitzung (F-536). *Form der Ablage*: **entfallen** vor dem Code (Prüfung des
  Architect). *Handbuch im Zielstand*: **eingetreten** (V-125), in diesem Slice für `record`
  behoben, die Rückkehr von `play:` trägt Folge-Slice `slice-v1-abschluss-einspielen`.
  *Name ohne Auflösen zwischen den Closures*: **entfallen**, DoD-Punkt 1 ist geliefert und am
  Binary bestätigt. Begründungen in §6. Die Randformen in §6 sind Entscheidungen, keine
  Risiken.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt kein Feld `liegt in`; nichts zu
  prüfen. Folge-Slice: `slice-v1-abschluss-einspielen` liegt in `next/`; `grep -n
  "slice-v1-abschluss-upstream-verbinden"` findet die Kennung in seinem §1 unter *Übernommen
  aus* (Wirkung bei `play`, U8, F-533, V-125) und in §3 (Zeile
  `docs/user/benutzerhandbuch.md`); sein §1 *Ausdrücklich NICHT* schließt weder `play` noch
  das Handbuch aus. `slice-v1-abschluss-tls-client` liegt in `next/`; die Kennung steht in
  seinem §1 unter *Übernommen aus*, mit `LH-FA-17.a` im Kopf; sein §1 *Ausdrücklich NICHT*
  schließt nur TLS zum Upstream und Client-Zertifikate aus, seine DoD führt die Kombinationen.
  Keiner liegt in `done/`. Register: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` tragen
  `evidence/slice-v1-abschluss-upstream-verbinden.md`; die übrigen genannten Einträge bestehen
  als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft die
  Paarungen erneut.

### Belege des Implementers

**Stand.** Geliefert sind DoD-Punkt 1 bis 3 in fünf Commits, jeder mit `make test` grün:
`fe849ac` (DoD-Punkt 1: Prüfung von `--upstream` und `PGWIRE_RECORDER_UPSTREAM` nach der
Zusammenführung, Namen nur aus der gewählten Datei, Hilfe), `694eb99` (DoD-Punkt 1 und 2:
benutzte Verbindung, `sslmode=require` bei `record`, Einsetzen, `PGR-E2005`, Port nach dem
Einsetzen, Zusammensetzen), `e35fb56` (DoD-Punkt 1: Integrationstests), `8ba8caa`
(DoD-Punkt 3: Handbuch), `80cb896` (Testfall zu einem grünen Mutanten, unten). Neu sind
`internal/adapters/driving/cli/upstream.go`, `internal/adapters/driving/cli/upstream_test.go`,
`internal/adapters/driving/cli/einsetzen_test.go` und
`test/integration/verbindung_e2e_test.go`; geändert `cli.go` (Feld `zuletzt` am Leser, Hilfe),
`datei.go` (`hatName`, `verbindung`; der Schlüssel `upstream` prüft über dieselbe Regel
`upstreamGueltig`), `verbindung.go` (`einsetzen`, `adresseRecord`, `zusammensetzen`),
`export_test.go` (`SetzeEin`), `leser_test.go` (`basis` gibt `--upstream` die Form
`host:port`), `verbindung_test.go` (`TestDateiUpstream` erwartet für den Namen die Adresse).
`internal/bootstrap` und `postgres.Upstream` sind unverändert; `RecordOptions.Upstream` trägt
die zusammengesetzte Adresse. Die Abdeckungstabellen sind in jedem Commit mit neuen
Deklarationen über `make abdeckung` nachgezogen.

**Randformen.** Keine Randform außerhalb von §6 entschieden, keine Rückgabe an den Architect;
kein Code-Commit ändert §6. Die Meldungstexte (`--upstream: weder Name einer Verbindung der
Konfigurationsdatei noch host:port`, `Umgebungsvariable PGWIRE_RECORDER_UPSTREAM: …`,
`connections.<Name>: sslmode=require ist bei record ungültig, …`, `connections.<Name>:
Umgebungsvariable <VAR> eines Platzhalters nicht gesetzt`, `connections.<Name>: Port nach dem
Einsetzen ist keine Zahl von 1 bis 65535`) nennen die Stelle nach U3 und keinen Wert, auch
die zur Kommandozeile nicht. Kein anderer Slice ist als neue Adresse genannt (§3.13).

**Größe.** Der Diff `9a04b12..80cb896` umfasst 712 hinzugefügte und 24 entfernte Zeilen: Code
141, Unit-Tests 424, Integration 94, Handbuch 34, Abdeckung 19; innerhalb der Schätzung von
660 bis 870 Zeilen (§6 *Risiken*, erster Punkt). Ob er in eine Review-Sitzung passt, urteilt
das Review.

**Läufe.**

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make test` | vor `fe849ac`, `694eb99`, `e35fb56`, `80cb896` | Exit 0 |
| `make lint` | vor `fe849ac`, `694eb99`, `e35fb56`, `80cb896` | Exit 0, `0 issues.` |
| `make test-integration` | vor `e35fb56` | Exit 0, `TestE2ERecordVerbindungAusDatei` und `TestE2ERecordVerbindungIPv6` PASS |
| `make docs-check` | vor `8ba8caa` | 0 Befunde |
| `make abdeckung` | je Commit mit Deklarationen | Tabellen nachgezogen; `make abdeckung-check` vor `8ba8caa` Exit 0 |
| `make gates` | `80cb896` | Exit 0 (darin `run-integration-tests: gruen`, `d-check: … 0 Befund(e)`, `lint-gegenprobe: gruen`) |
| `make gates` | Commit dieser Belege | Exit 0, vor der Übergabe |

**Weg der Mutanten.** Je Mutant ein frisches Verzeichnis unter dem Scratch-Pfad der Sitzung,
in das `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`, `.golangci.yml`, `Dockerfile` und
`.dockerignore` mit `cp -r` ohne `-p` kopiert werden (neue mtime); ein Skript ersetzt genau
ein Vorkommen und bricht bei einer anderen Trefferzahl ab. Unit-Mutanten: `docker build
--target test` in der Kopie (gofmt, vet, `go test ./...`), Integrations-Mutanten:
`tools/test/run-integration-tests.sh` in der Kopie, mit eigenem Netz, Container und Volume,
die das Skript danach entfernt. Danach werden das Verzeichnis und das Image des Mutanten
gelöscht; der Arbeitsbaum blieb unberührt. Gefahren: 34 Unit-Mutanten und 4
Integrations-Mutanten; 37 rot, einer zuerst grün (unten), mit dem Test aus `80cb896` rot. Ein
Mutant (*config show setzt ein*) endete zuerst mit einer Nil-Dereferenz in einem fremden
Test, ist mit Prüfung auf eine fehlende Datei neu gefahren und dann aus dem richtigen Grund
rot.

**DoD-Punkt 1 — Wert von `--upstream`, benutzte Verbindung, Adresse.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Wert der Kommandozeile geprüft (A7, A8) | Prüfung der Kommandozeile entfernt | `TestUpstreamUngueltig`, `TestUpstreamNameNurAusDatei` |
| Wert der Umgebungsvariable geprüft (A7, A8) | Prüfung der Umgebung entfernt | `TestUpstreamUngueltig`, `TestUpstreamNameNurAusDatei` |
| auch der Wert, der nicht gilt (A8) | Umgebung nur ohne Kommandozeile geprüft | `TestUpstreamUngueltig` |
| Kommandozeile vor Umgebungsvariable (*Fehler*) | Reihenfolge getauscht | `TestUpstreamUngueltig` |
| Meldung zur Umgebung ohne Wert | Wert in die Meldung | `TestUpstreamUngueltig` |
| Host nach F-521 bei Option und Umgebung | Prüfung nur am Schlüssel (beide Prüfungen oben entfernt) | `TestUpstreamUngueltig` (`GEHEIM@h:5`, `GEHEIM h:5`, `a/GEHEIM:5`, `h%41GEHEIM:1`, `a[GEHEIM:1`, Steuerzeichen) |
| Namen nur aus der gewählten Datei (U1) | Namen gegen eine leere Menge | `TestUpstreamNameNurAusDatei`, `TestDateiUpstream`, `TestDateiUpstreamUngueltig` |
| Name nur bei genauer Übereinstimmung | auch in Kleinbuchstaben angenommen | `TestUpstreamNameNurAusDatei`, `TestDateiUpstreamUngueltig` |
| Prüfung nach Pflichtoptionen (*Fehler*) | Prüfung in der Schleife der Zusammenführung | `TestUpstreamReihenfolge` |
| Hilfe nennt host:port oder Namen (U7) | alter Text | `TestUpstreamHilfe` |
| benutzt ist die Verbindung, die der Wert nennt | erste Verbindung der Datei | `TestUpstreamVerbindung`, `TestUpstreamVariableFehlt`, `TestUpstreamReihenfolgeAmEnde` |
| der zusammengeführte Wert, nicht die Umgebung | Name aus der Umgebung, wenn gesetzt | `TestUpstreamVerbindung` |
| auch aus dem Schlüssel `upstream` | nur aufgelöst, wenn Option oder Umgebung gesetzt | `TestUpstreamVerbindung`, `TestDateiUpstream` |
| `record` verbindet zu Host und Port der Verbindung | Adresse nicht übernommen | `TestUpstreamVerbindung` und weitere 5 Unit-Tests; E2E: `TestE2ERecordVerbindungAusDatei`, `TestE2ERecordVerbindungIPv6` |
| Host mit `:` in eckigen Klammern (U5, F-528) | nie geklammert | `TestUpstreamVerbindung`, `TestUpstreamEinsetzen`; E2E: `TestE2ERecordVerbindungIPv6` |
| Host ohne `:` ohne Klammern (U5) | immer geklammert | `TestUpstreamVerbindung`, `TestUpstreamEinsetzen` und 4 weitere |
| Port wie eingesetzt (U4) | Port aus dem Teil vor dem Einsetzen | `TestUpstreamEinsetzen`; E2E: `TestE2ERecordVerbindungAusDatei` |
| Zeile beim Start nennt die Adresse, nicht den Namen (U6) | Adresse nicht übernommen (die Zeile nennt den geschriebenen Wert) | E2E: `TestE2ERecordVerbindungAusDatei` (`upstream=` mit Host und Port, ohne `staging`) |
| bei `record` bleiben Benutzer, Passwort, Datenbank unbeachtet | Einsetzen auch in diese Teile | `TestUpstreamVerbindung`; E2E: `TestE2ERecordVerbindungAusDatei` |

**DoD-Punkt 2 — Einsetzen, `PGR-E2005`, Port, `sslmode`, Reihenfolge.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| nicht gesetzte Variable ist `PGR-E2005` | Code `PGR-E2004` statt `PGR-E2005` | `TestUpstreamVariableFehlt`, `TestEinsetzenAlleTeile`, `TestUpstreamReihenfolgeAmEnde` |
| leere Variable gilt als nicht gesetzt | leerer Wert eingesetzt | `TestUpstreamVariableFehlt`, `TestEinsetzenAlleTeile`, `TestUpstreamReihenfolgeAmEnde` |
| Meldung nennt die erste Variable in der Reihenfolge der URL (A12) | Host und Port vertauscht eingesetzt | `TestUpstreamVariableFehlt` |
| innerhalb eines Teils die erste (A12) | Name der letzten Variable des Teils | `TestUpstreamVariableFehlt`, `TestEinsetzenAlleTeile` |
| Meldung nennt `connections.<Name>` (U3) | Stelle `connections` ohne Name | `TestUpstreamPortNachEinsetzen`, `TestUpstreamReihenfolgeAmEnde` |
| Variablen nicht benutzter Verbindungen unbeachtet | jede Verbindung eingesetzt | `TestUpstreamVariableFehlt`, `TestUpstreamReihenfolgeAmEnde` |
| `config show` meldet kein `PGR-E2005` | `config show` setzt die Verbindung des Schlüssels `upstream` ein | `TestUpstreamVariableFehlt` |
| einmal eingesetzt, keine erneute Auswertung | eingesetzter Wert erneut ersetzt | `TestUpstreamEinsetzen` |
| Wert unverändert, nicht dekodiert | eingesetzter Wert dekodiert | `TestUpstreamEinsetzen`, `TestEinsetzenAlleTeile` |
| Wert unverändert, nicht gekürzt | Leerraum am Rand entfernt | `TestUpstreamEinsetzen`, `TestUpstreamPortNachEinsetzen` (seit `80cb896`, unten) |
| Port nach dem Einsetzen geprüft, `PGR-E2004` (Rückgabe 8) | Prüfung entfernt | `TestUpstreamPortNachEinsetzen` |
| Meldung zum Port ohne Wert (U3) | Wert in die Meldung | `TestUpstreamPortNachEinsetzen` (`GEHEIM`) |
| den Host prüft der Start nicht (Rückgabe 8) | eingesetzter Host wie ein wörtlicher geprüft | `TestUpstreamEinsetzen` (`a:b`, `a@b/%41`) |
| `sslmode=require` bei `record` `PGR-E2004` | Prüfung entfernt | `TestUpstreamReihenfolgeAmEnde` |
| `sslmode=require` vor den Variablen (A9) | nach dem Einsetzen geprüft | `TestUpstreamReihenfolgeAmEnde` |
| `--upstream` vor `sslmode=require` | `sslmode=require` vor der Prüfung des Werts | `TestUpstreamReihenfolgeAmEnde` |
| Variablen vor dem Port | Form des Ports vor dem Einsetzen der Variablen geprüft | `TestUpstreamVariableFehlt` (leere Variable im Port ist `PGR-E2005`), `TestUpstreamReihenfolgeAmEnde`, `TestUpstreamEinsetzen` |
| Einsetzen gilt für jeden Teil, erste fehlende in der Reihenfolge der URL (U8) | Name der letzten Variable des Teils; leerer Wert eingesetzt; Wert dekodiert (Zeilen oben) | `TestEinsetzenAlleTeile` über `SetzeEin`, je Teil Benutzer, Passwort, Host, Port, Datenbank |
| Variablen einmal gelesen (U2), erste Hälfte: die Adresse steht nach `Parse` eingesetzt in `RecordOptions.Upstream` | Adresse nicht in `c.Record.Upstream` übernommen (Mutant A1 der Verifikation; dieselbe Mutation wie *Adresse nicht übernommen* in DoD-Punkt 1) | `TestUpstreamEinmalGelesen` (`einsetzen_test.go:251`: Name `v` statt `erst:5432`), dazu `TestUpstreamVerbindung`, `TestUpstreamEinsetzen`, `TestUpstreamVariableFehlt`, `TestUpstreamReihenfolgeAmEnde`, `TestDateiUpstream` |
| Variablen einmal gelesen (U2), zweite Hälfte: ein `Setenv` nach `Parse` wirkt nicht | Grenze, unten | — |

**Grüne Mutanten, eingeordnet.**

- *Eingesetzter Wert am Rand gekürzt* (`strings.TrimSpace`) blieb grün: Er ändert das
  Verhalten (Host mit Leerraum, Port `" 5"`), ist über die Schnittstelle fangbar, und kein
  Fall hatte Leerraum am Rand. `80cb896` ergänzt in `TestUpstreamEinsetzen` den Host `" h\t"`
  und in `TestUpstreamPortNachEinsetzen` die Ports `" 5"` und `"5\n"`; danach rot in beiden.
- Sonst kein grüner Mutant.

**Grenzen.**

- *U2, Einsetzen beim Verbinden statt beim Lesen* ist im Zuschnitt nicht baubar, ohne
  `postgres.Upstream` zu ändern (§1, Schicht-Abgrenzung; §6 U2). Die zweite Hälfte von
  `TestUpstreamEinmalGelesen` (`Setenv` nach `Parse`) kann deshalb nicht rot werden; sie
  sichert nur die Form der Übergabe: `RecordOptions.Upstream` ist ein fester Text. Die erste
  Hälfte hält Mutant A1 (Tabelle oben).
- *Vorrang des Namens vor `host:port`* (A7) hat keinen Fall: Ein Name enthält kein `:`,
  `host:port` immer eines (akzeptiertes Negativ in §6, A7).
- Die einzelnen Zeichen der Regel F-521 prüft `hostPortForm`, mit Mutanten je Zeichen belegt im
  Lade-Slice am Schlüssel `upstream`; hier belegt die Mutation *Prüfung nur am Schlüssel* (die
  zwei Prüfungen von Option und Umgebung entfernt), dass Option und Umgebung dieselbe Prüfung
  durchlaufen, mit einem Fall je Zeichenklasse in `upstreamUngueltig`.

**Nacharbeit zum Review** (`docs/reviews/2026-10-09-review-slice-v1-abschluss-upstream-verbinden.md`,
Stand `a482303`).

- *F-534, Auswahl der Verbindung nach dem genauen Namen.* Merkmal: Die Suche in
  `datei.verbindung` vergleicht den Namen; `hatName` davor hält die Auswahl nicht.
  `TestUpstreamVerbindung` führt jetzt eine Datei mit `Staging`, `stag`, `staging2` und
  `staging` in dieser Reihenfolge und erwartet je Name genau die eigene Adresse, über
  `--upstream`, `PGWIRE_RECORDER_UPSTREAM` und den Schlüssel `upstream` (die Sonde aus dem
  Review: `--upstream=staging` ergibt `klein:5432`). Neben dem gemeldeten Mutanten je weitere
  Ausprägung des Merkmals eine Mutation:

  | Zusage | Mutation | rote Tests |
  |---|---|---|
  | Auswahl genau in Groß- und Kleinschreibung | Mutant S, `strings.EqualFold(v.name, name)` | `TestUpstreamVerbindung` |
  | ebenso | beide Namen kleingeschrieben verglichen | `TestUpstreamVerbindung` |
  | kein Präfix des Werts trifft | `strings.HasPrefix(name, v.name)` | `TestUpstreamVerbindung` (`stag` vor `staging`) |
  | kein Name, der mit dem Wert beginnt, trifft | `strings.HasPrefix(v.name, name)` | `TestUpstreamVerbindung` (`staging2` vor `staging`) |

  Zuerst blieb der dritte Mutant grün, weil `stag` hinter `staging` stand; mit der Reihenfolge
  oben rot. Weg wie oben (frische Kopie mit `cp -r` ohne `-p`, `docker build --target test`).
- *F-530.* Der Kommentar an `SetzeEin` nennt jetzt `""` für Benutzer, Passwort und Datenbank,
  die die URL nicht schreibt, und `5432` für einen nicht geschriebenen Port; so prüft es
  `TestEinsetzenAlleTeile` (Fall *ohne Benutzer*).
- *F-531.* Die Kommentare an `option.zuletzt` und `lies` nennen, dass `zuletzt` den Wert ersetzen
  darf, dass `upstreamRecord` `c.Record.Upstream` auf die Adresse setzt und dass ein späterer
  Aufruf den ersetzten Wert sieht.
- *F-533.* Aus dem Handbuch §5 *Konfigurationsdatei* sind die Sätze zu `play` heraus: das
  Passwort aus dem Platzhalter (der Absatz beginnt jetzt mit dem Einsetzen des Platzhalters),
  *„Beim Einspielen verbindet `require` …“* mit `--upstream-tls` und *„gilt
  `PGWIRE_RECORDER_PASSWORD`“*. `slice-v1-abschluss-einspielen` nahm sie nicht an (§1 und DoD
  ohne Handbuch, kein Ausschluss); dort steht die Sendung jetzt in §1 unter *Übernommen aus*
  diesem Slice und in §3, mit der Kennung dieses Slice; hier in §1 *Ausdrücklich NICHT* und §3.
  Stehen geblieben sind Stellen außerhalb der Absätze aus DoD-Punkt 3, die `play` schon
  vorher nannten (Einleitung von §5, Zeile `PGR-E4005` in §7); das Beispiel mit `play:` gehört
  in DoD-Punkt 3 und ist mit der Nacharbeit zur Verifikation umgestellt (V-125, unten).
- *F-532, F-535* sind Hinweise an Architect und Verifier; kein Code.

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make test`, `make lint` | Nacharbeit, vor dem Commit | Exit 0, `0 issues.` |
| `make docs-check`, `make kopf-check` | Nacharbeit, vor dem Commit | 0 Befunde, Exit 0 |
| `make gates` | Commit der Nacharbeit | Exit 0, vor der Übergabe |

**Nacharbeit zur Verifikation** (`docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-upstream-verbinden.md`,
Stand `e74077b`).

- *V-125, Beispiel in §5 Konfigurationsdatei.* Merkmal: Ein Text im Handbuch zeigt eine
  Eingabe, die das gelieferte Binary ablehnt. Ausprägungen in §5 *Konfigurationsdatei*: (1) der
  Abschnitt `play:` im Beispiel, (2) `play:` in der Abschnittsliste davor, (3) `--upstream
  staging` und `PGWIRE_RECORDER_UPSTREAM=staging` im Satz nach dem Beispiel: Die einzige
  Betriebsart mit `--upstream` ist heute `record`, und dort ist `staging` wegen
  `sslmode=require` `PGR-E2004`. Das Beispiel zeigt jetzt einen Abschnitt `record:` mit
  `upstream: lokal`, `listen: 127.0.0.1:15432` und `output: ./recordings/users.yaml`, die
  Liste nennt `record:` und `replay:`, der Satz nennt `lokal`. Die Verbindung `staging` bleibt
  im Beispiel, weil sie Platzhalter und `sslmode` zeigt; der Absatz zu `record` sagt, dass
  `require` dort ungültig ist. Die Rückkehr von `play:` geht an `slice-v1-abschluss-einspielen`
  (dort §1 *Übernommen aus* diesem Slice und §3, mit der Kennung dieses Slice; sein §1
  *Ausdrücklich NICHT* schließt das Handbuch nicht aus, die Zeitangaben gehen an
  `slice-v1-abschluss-zeitangaben`, deshalb nennt die Sendung `keep_timing` und `timing_mode`
  nicht).

  Probe mit dem gebauten Image (`make build`, `pgwire-recorder:dev`), das Beispiel per `awk`
  wörtlich aus dem Handbuch in eine Datei unter dem Scratch-Pfad der Sitzung, Lauf mit
  `docker run --network none`:

  | Probe | Stand des Handbuchs | Ergebnis |
  |---|---|---|
  | `config show --config beispiel.yaml` | vor der Nacharbeit | Exit 2, `Konfiguration [PGR-E2004]: Konfigurationsdatei: play: unbekannter Schlüssel` |
  | `config show --config beispiel.yaml` | nach der Nacharbeit | Exit 0, erste Zeile der Pfad, danach der Inhalt der Datei |
  | `record --config beispiel.yaml`, nach 2 s `docker stop` | nach der Nacharbeit | Exit 0, `record gestartet listen=127.0.0.1:15432 upstream=localhost:5432`, `Herunterfahren begonnen sessions=0`, `record beendet output=./recordings/users.yaml` |
  | `record --config beispiel.yaml --upstream staging` | nach der Nacharbeit | Exit 2, `connections.staging: sslmode=require ist bei record ungültig, …` (Grund für Ausprägung 3) |

- *V-126, U2.* Die Tabelle zu DoD-Punkt 2 führt U2 in zwei Zeilen: die erste Hälfte von
  `TestUpstreamEinmalGelesen` mit Mutant A1 (`c.Record.Upstream = adresse` in `upstreamRecord`
  durch `_ = adresse` ersetzt), die zweite als Grenze. A1 selbst gefahren: frische Kopie wie in
  *Weg der Mutanten* oben, `docker build --target test`, rot in sechs Tests, darunter
  `TestUpstreamEinmalGelesen` an `einsetzen_test.go:251`; Kopie und Image danach gelöscht.

Kein Produkt-Code und kein Test geändert; keine Randform entschieden.

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make build` | Nacharbeit, vor dem Commit | Exit 0 |
| `make docs-check`, `make kopf-check` | Nacharbeit, vor dem Commit | `d-check: 405 Datei(en) geprüft, 0 Befund(e)`, Exit 0 |
| `make gates` | Commit der Nacharbeit | Exit 0, vor der Übergabe |

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
