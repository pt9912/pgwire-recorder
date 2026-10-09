# Slice slice-v1-abschluss-verbindungen-platzhalter: Benannte Verbindungen und Platzhalter in der Konfigurationsdatei

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

**Ziel:** Die benannten Verbindungen der Konfigurationsdatei (`connections`) werden nach der Grammatik der URL samt `sslmode` geprüft, `--upstream` und der Schlüssel `upstream` lösen den Namen einer Verbindung auf, `${VAR}` wird aus der Umgebung eingesetzt und `$$` steht für `$`; eine nicht gesetzte Variable der benutzten Verbindung ist `PGR-E2005`, ein Klartext-Passwort `PGR-E2006`.

**Übernimmt:** `slice-v1-abschluss-konfigurationsdatei` — dessen DoD-Punkt 3 (*Verbindungen
und Platzhalter*) nach dem vorab benannten Schnitt aus seinem §4 (*zu groß*, eingetreten am
2026-10-09; dort §1, *Abgegeben*). Dorthin kam der Punkt aus `slice-v1-abschluss-konfiguration`
und davor aus `slice-v1-abschluss-betrieb`. Im Einzelnen:

- die Grammatik der URL einer Verbindung und `sslmode`, auch `sslmode=require` bei `record`
  (`PGR-E2004`), in jeder Verbindung der Datei, auch einer nicht benutzten;
- `${VAR}` und `$$` in jedem Wert der Datei, ein Platzhalter außerhalb einer URL
  (`PGR-E2004`) und der Port nach dem Einsetzen;
- `PGR-E2005` und `PGR-E2006` — die Konstanten liegen seit jenem Slice ohne Erzeuger in
  der Code-Tabelle, dieser Slice erzeugt sie;
- das Auflösen eines Verbindungsnamens in `--upstream` und im Schlüssel `upstream`; bei
  `record` gelten Host und Port;
- der Teil des Benutzerhandbuchs dazu: in §5 *Konfigurationsdatei* die Absätze zu
  Verbindungen, `sslmode` und Platzhaltern, in §7 *Fehlercodes* die Zeilen `PGR-E2005` und
  `PGR-E2006`;
- die Randformen zu Verbindungen und Platzhaltern aus §6 jenes Slice, dazu die
  Entscheidungen des Architect vom 2026-10-08 dazu (Rückgaben 5, 6, 7, 8 und 10, §6 unten).

**Aufsetzen.** Der Slice setzt auf dem Laden von `slice-v1-abschluss-konfigurationsdatei` auf:
Dort ist `connections:` eine Abbildung, ein Name einer Verbindung hat die Form eines Werts
ohne Steuerzeichen, und ein Wert ist ein Skalar (dort Rückgaben 4 und 11); bis zu diesem
Slice gilt ein Wert mit `$` wörtlich, und eine URL wird nicht zerlegt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Wahl, Laden und Form der Datei, die Priorität über drei Quellen, `config show`, die Hilfe
  und die Konstanten `PGR-E2004` bis `PGR-E2006` — `slice-v1-abschluss-konfigurationsdatei`;
  er liegt vor diesem Slice, und dieser setzt auf ihm auf.
- Die Wirkung einer Verbindung bei `play` (Benutzer und Datenbank der URL mit dem Vorrang von
  `--user` und `--database`, das Passwort aus dem eingesetzten Platzhalter, TLS nach
  `sslmode`) — `slice-v1-abschluss-einspielen` (dort §1, *Übernommen aus* diesem Slice);
  `play` gibt es erst mit ihm. Hier wird `sslmode` geprüft, und das Auflösen des Namens ist
  an `record` getestet.
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

- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Jede Verbindung der Datei, auch eine nicht benutzte, hat die Form der
      URL aus `LH-FA-17.a` (Schema, Host, Port mit Default `5432` und Wert 1 bis 65535,
      Datenbank, Parameter, Prozent-Dekodierung wörtlicher Teile), der einzige Parameter
      `sslmode` nimmt `disable` und `require`, jede andere Form ist `PGR-E2004` mit dem
      Namen der Verbindung; `sslmode=require` der benutzten Verbindung ist bei `record`
      `PGR-E2004`. `--upstream` und der Schlüssel `upstream` lösen einen Namen nur bei
      genauer Übereinstimmung auf, mit Vorrang vor `host:port`; ein Wert, der weder Name
      noch `host:port` ist, ist `PGR-E2001`, aus der Datei `PGR-E2004`; `record` verbindet
      zu Host und Port der Verbindung (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `${VAR}` wird für die benutzte Verbindung einmal in die zerlegte URL
      eingesetzt, der Wert steht unverändert in seinem Teil, nur der Port wird danach
      geprüft (`PGR-E2004`); die Form des Platzhalters prüft das Laden in jedem Teil jeder
      Verbindung, ein Platzhalter außerhalb einer URL ist `PGR-E2004`, und `$$` steht in jedem
      Wert der Datei für `$`. Eine nicht gesetzte oder leere Variable der benutzten Verbindung
      ist `PGR-E2005`; die Variablen nicht benutzter Verbindungen und bei `record` die in
      Benutzer, Passwort und Datenbank bleiben unbeachtet; `config show` meldet kein
      `PGR-E2005` und zeigt Platzhalter unaufgelöst (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Ein Klartext-Passwort ist `PGR-E2006` in jeder Verbindung der Datei, auch
      bei `config show` — ein Passwortteil, der nicht genau ein `${VAR}` ist, auch ein
      fehlerhafter Platzhalter dort, und ein Parameter `password`; die Prüfung folgt der
      Reihenfolge aus `LH-FA-17.a` innerhalb einer URL und am Ende (`--upstream`, `PGR-E2005`,
      Port nach dem Einsetzen) (Test). Das Benutzerhandbuch beschreibt Verbindungen,
      `sslmode`, Platzhalter, `PGR-E2005` und `PGR-E2006` wie geliefert. Beleg in §7 für
      alle drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` (`LH-FA-17.a`) | keine Änderung geplant | Die Randformen in §6 sind entschieden (Architect, 2026-10-08, vor dem Code); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driving/cli` | update | Zerlegen der URL mit Platzhaltern, Grammatik, Prozent-Dekodierung, `sslmode`; `$$` in jedem Wert der Datei, Platzhalter außerhalb einer URL; Klartext-Passwort; Auflösen des Namens bei `--upstream` und `upstream`; Einsetzen der Variablen der benutzten Verbindung und Prüfung des Ports danach; `PGR-E2005` und `PGR-E2006` |
| `internal/bootstrap` | update, falls nötig | nur, wenn die aufgelöste Adresse anders als heute an `record` übergeben werden muss |
| `internal/adapters/driving/cli` (Unit-Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a`, je Randform aus §6 ein Fall |
| `docs/user/benutzerhandbuch.md` | update, falls abweichend | §5 *Konfigurationsdatei* (Verbindungen, `sslmode`, Platzhalter) und §7 *Fehlercodes* (`PGR-E2005`, `PGR-E2006`) beschreiben den Zielstand; nachgezogen wird, was der gelieferte Stand anders sagt |
| `docs/user/abdeckung-*.md` | update | über `make abdeckung` aus den Deklarationen der neuen Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfigurationsdatei` liegt in `done/`
(Laden der Datei, Form von `connections:`, Konstanten `PGR-E2005` und `PGR-E2006`). Schritt 5
der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (Schnitt vom
2026-10-09). Die Randformen aus §6 entschied der Architect am 2026-10-08 vor dem Code; vor
dem ersten Code-Commit prüft er die Liste noch einmal gegen diesen Zuschnitt
(`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Kern wird nötig. Schnitt dann: *URL und
  Auflösen* (DoD-Punkt 1) und *Platzhalter und Passwort* (DoD-Punkt 2 und 3); der zweite
  setzt den ersten voraus, weil Platzhalter in der zerlegten URL eingesetzt werden.
- `in-progress` → `open` (blockiert — Carveout?): Eine Zusage aus `LH-FA-17.a` zu Verbindungen
  oder Platzhaltern widerspricht dem Laden von `slice-v1-abschluss-konfigurationsdatei` (etwa
  der Form eines Werts) und verlangt eine neue Entscheidung; dann zuerst die Entscheidung.

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
Schnitt vom 2026-10-09 aus §6 von `slice-v1-abschluss-konfigurationsdatei` hierher; dorthin
kamen die ersten aus §6 von `slice-v1-abschluss-konfiguration`, geprüft vom Architect am
2026-10-08 vor dem Code. Die Rückgaben unten entschied er am 2026-10-08 vor dem Code, die
Nummern sind die aus jenem Plan. Offen ist keine.

*Verbindungen und Platzhalter*

- **Name bei `--upstream`** — nur bei genauer Übereinstimmung; ein Wert, der weder Name noch
  `host:port` ist, ist ungültig (`PGR-E2001`, aus der Datei `PGR-E2004`); `LH-FA-17.a`.
- **`sslmode=require` bei `record`** — `PGR-E2004` der benutzten Verbindung; `LH-FA-17.a`.
- **Klartext-Passwort** — `PGR-E2006`, in jeder Verbindung der Datei; `LH-FA-17.a`.
- **Variable eines Platzhalters der benutzten Verbindung nicht gesetzt** — `PGR-E2005`;
  `LH-FA-17.a`.
- **Variable eines Platzhalters gesetzt, aber leer** — nicht gesetzt, `PGR-E2005`;
  Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Platzhalter in Teilen der URL, die `record` ignoriert** (Benutzer, Passwort, Datenbank) —
  unbeachtet; `record` löst nur Platzhalter in Host und Port auf; Entscheidung des Nutzers
  vom 2026-10-08, `LH-FA-17.a`.
- **`$${VAR}`, Rekursion, Form des Namens, Einsetzen in die URL** — `$$` steht für `$` in
  jedem Wert der Datei; Name `[A-Za-z_][A-Za-z0-9_]*`, sonst `PGR-E2004`; einmal eingesetzt;
  URL vor dem Einsetzen zerlegt, Wert unverändert in seinem Teil; Entscheidung des Nutzers
  vom 2026-10-08, `LH-FA-17.a`. Hinweis an den Implementer: Die Zerlegung der
  Standardbibliothek lehnt einen Port `${PORT}` ab; zerlegt wird mit Platzhaltern.
- **Platzhalter außerhalb einer URL** — `PGR-E2004` (Tabelle *Fehler*); `LH-FA-17.a`.

*Anzeige und Fehler, soweit sie Verbindungen betreffen*

- **`config show` und `PGR-E2005`** — kommt nicht vor; `config show` setzt keine Variable ein
  und zeigt Platzhalter unaufgelöst; `LH-FA-17.a`.
- **`config show` und Geheimnisse** — keine Maskierung, und keine nötig: Ein Passwort steht in
  keiner gültigen Datei (`PGR-E2006`, mit diesem Slice), Platzhalter erscheinen unaufgelöst;
  `LH-FA-17` und `LH-FA-17.a` (*Anzeige*). Die Randform selbst bleibt in
  `slice-v1-abschluss-konfigurationsdatei`; hier kommt ihre Voraussetzung.
- **Reihenfolge, Teil der Verbindungen** — innerhalb einer URL in der Reihenfolge ihrer Teile;
  am Ende nach der Zusammenführung `--upstream`, die Variablen der Platzhalter der benutzten
  Verbindung (`PGR-E2005`) und danach ihr Port nach dem Einsetzen; `LH-FA-17.a`. Den übrigen
  Teil der Reihenfolge liefert `slice-v1-abschluss-konfigurationsdatei`.

*Rückgaben des Implementers vom 2026-10-08*, entschieden vom Architect am 2026-10-08 vor dem
Code (`f19b440`), alle in `LH-FA-17.a`, abgeleitet aus den Grundsätzen dort: Text des Skalars
zählt, Laden ohne Kommando, geschlossene Fehlertabelle, URL vor dem Einsetzen zerlegt, Abbruch
beim ersten Fehler. Die Rückgaben 1 bis 4, 9, 11 und 12 liefert
`slice-v1-abschluss-konfigurationsdatei`.

5. **Nicht benutzte Verbindung** (Teil der Verbindungen aus Rückgabe 5) — mitgeprüft, auch
   Grammatik, `sslmode` und Klartext-Passwort; vom Kommando hängen nur `sslmode=require` bei
   `record` und die Variablen der Platzhalter ab. Dass der Abschnitt eines anderen Kommandos
   mitgeprüft wird, liefert `slice-v1-abschluss-konfigurationsdatei`.
6. **URL ohne Port** — Port `5432`; die Grammatik lässt den Port weg. Ein Port sind Ziffern
   mit Wert 1 bis 65535, sonst `PGR-E2004`.
7. **Platzhalter-Syntax in ignorierten Teilen** — geprüft beim Laden, `PGR-E2004`;
   unbeachtet ist dort nur die Variable. Im Passwortteil ist ein fehlerhafter Platzhalter
   ein Klartext-Passwort (`PGR-E2006`).
8. **Eingesetzter Wert macht seinen Teil ungültig** — nur der Port wird nach dem Einsetzen
   geprüft: `PGR-E2004`, im letzten Schritt direkt nach den Variablen der benutzten
   Verbindung (`PGR-E2005`). Den Host prüft der Start nicht; ein Host mit `/` scheitert beim
   Verbindungsaufbau.
10. **URL-Sonderformen** — was die Grammatik nicht zulässt, ist `PGR-E2004`: Schema
    `postgres://`, leerer Host, fehlende oder leere Datenbank, Fragment, Parameter ohne `=`,
    ein Parameter zweimal (auch `sslmode`). Wörtliche Teile werden prozent-dekodiert, ein
    ungültiges Escape ist `PGR-E2004`; Platzhalter und `$$` gelten vor der Dekodierung, ein
    eingesetzter Wert wird nicht dekodiert. Innerhalb einer URL gilt die Reihenfolge ihrer
    Teile.

**Risiken:**

- Die Zerlegung der URL mit Platzhaltern braucht einen eigenen Zerleger, weil die
  Standardbibliothek `${PORT}` als Port ablehnt (Hinweis oben); Grammatik, Dekodierung und
  Einsetzen in einem Zerleger können den Diff über eine Review-Sitzung heben. Gegenmittel:
  der vorab benannte Schnitt in §4 — **Ausgang:** offen bis Closure.
- Das Benutzerhandbuch beschreibt Verbindungen und Platzhalter schon im Zielstand und kann vom
  gelieferten Stand abweichen (§3) — **Ausgang:** offen bis Closure.
- Bis zur Closure dieses Slice lehnt der gelieferte Stand ein Klartext-Passwort in der Datei
  nicht ab, und `config show` zeigt es (Risiko in §6 von
  `slice-v1-abschluss-konfigurationsdatei`); DoD-Punkt 3 schließt die Lücke —
  **Ausgang:** offen bis Closure.

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `cfc4d6b` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — dieser Slice entsteht aus dem Schnitt von
  `slice-v1-abschluss-konfigurationsdatei`, der durch Übernahmen über eine Review-Sitzung
  wuchs. Er selbst übernimmt nur dessen DoD-Punkt 3 und bekommt drei Liefer-Punkte daraus; er
  ist kein weiterer Beleg. Ob der Geber bei seiner Closure der dritte ist, entscheidet dessen
  Closure.
- `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (1×) — der Schnitt läuft diesmal formal über
  `next/` (Entscheidung des Nutzers vom 2026-10-09 zu V-117); kein weiterer Beleg.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (15×, verkörpert in `AGENTS.md` §3.12)
  und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (4×, verkörpert) — alle
  Randformen in §6 sind vor dem Code entschieden; eine weitere geht an den Architect zurück.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (16×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (22×, §3.11) — je Zusage eine Mutation,
  Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (18×, §3.9) — jede Korrektur zieht §1, §3 und §6 im
  selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — der Geber nennt diesen Slice in
  §1 unter *Abgegeben*, der Nehmer `slice-v1-abschluss-einspielen` nennt ihn in §1 unter
  *Übernommen aus*, beides im selben Commit.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor dem
Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
