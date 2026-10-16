# Slice slice-v1-abschluss-container: Container-Image und Betriebsdokumentation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031` · `SPEC-035` · [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (gelesen: das Gate hält die zwei Handbuch-Tabellen)

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

**Ziel:** Ein Docker/OCI-Image (`linux/amd64`, `linux/arm64`) und ein reproduzierbarer Binary-Build existieren, und die Betriebsdokumentation nennt Optionen, Exit-Codes, Meldungskatalog und verwendete Umgebungen.

**Übernommen aus `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen) und `slice-v1-abschluss-homebrew`** (dort je §1, Abgrenzung): das Container-Image (Docker/OCI); das ist das Ziel oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Veröffentlichungs-Mechanismus der Registries (`ghcr.io`, `docker.io`) — Gegenstand des Release-Verfahrens in `docs/maintainer/releasing.md`; dieser Slice liefert das Image.
- Gate für Code-Tabelle und Katalog — liefert `slice-harness-meldungskatalog-gate` (wellenlos, Entscheidung des Nutzers vom 2026-10-07) direkt vor diesem Slice (Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)); das Gate liegt in `done/` und prüft seit `slice-harness-meldungskatalog-gate` (`SPEC-054` Punkt 1, 4c, 4d) auch die Tabellen `### Fehlercodes` und `### Warnungen` des Handbuchs: Ihre Codes sind genau die der Code-Tabelle, andere Betriebsdokumente liest es nicht. Schreibt dieser Slice diese Tabellen fort, hält er `make meldungskatalog-check` grün; das Gate selbst hält die Gleichheit, ein eigener Liefer-Punkt dafür entfällt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Das Image startet `replay` ohne PostgreSQL und besteht den Walking-Skeleton-Smoke im Container. Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung). Aus `slice-doku-ist-stand` (dort §6, *Installation*): Der Slice liefert im Handbuch §2 *Container* und *Das Binary aus dem Image* (Image und seine Plattformen, statt „für die Architektur des Rechners, auf dem Sie gebaut haben“) und die Plattform-Zeile in §1 *Voraussetzungen* („Zum Bauen: … Zum Ausführen: Linux auf der Architektur des Rechners“), soweit das gebaute Image sie belegt.
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus): Das Binary ist für alle Zielplattformen aus `SPEC-035` reproduzierbar gebaut (zwei Builds desselben Quellstands, dieselbe Prüfsumme), das Image ist gebaut; Reproduzierbarkeit verlangt `LH-QA-01` nur für das Binary. Die Betriebsdokumentation ist vorhanden und benennt die Abnahme-Umgebungen.
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
| `Dockerfile`, `harness/mk/…` | neu/update | Image-Build per `make`, Docker-only |
| `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md` | update | auf den Stand der umgesetzten Funktionen prüfen und anpassen |
| `test/integration` | update | Smoke im Container |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die anderen Slices der Welle sind `done`, und `slice-harness-meldungskatalog-gate` liegt in `done/` (§1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- Dockerfile harness tools test docs/user docs/maintainer README.md` liegt über 1200; Basis ist der letzte Commit vor dem Anspruch. Marke und Schätzung (Planner, 2026-10-17, vorläufig; der Architect schärft sie im Vorlauf vor dem Code, `AGENTS.md` §3.14): Spanne 550 bis 1050, Mitte 800, Meldeschwelle 1000 (Dockerfile und Fragmente 80 bis 150, Prüfung der Reproduzierbarkeit je Zielplattform 60 bis 120 und ihre Gegenprobe 150 bis 300 nach dem Faktor Zwei bis Drei, Smoke im Container 100 bis 200, Handbuch, README und Betriebsdokumentation 150 bis 300). Geber: `slice-harness-meldungskatalog-gate` 607 (Skript mit Gegenprobe), `slice-v1-abschluss-verbindungsende-doku` 287 (Handbuch). Die Schätzung liegt unter der Marke, ein Schnitt vor dem Code ist nicht geboten.
- `in-progress` → `next`: die Doku wächst über einen Liefer-Punkt hinaus — zurück zur Zerlegung.
- `in-progress` → `open`: Das Binary lässt sich ohne Entscheidung zu Build-Flags und Toolchain-Pin nicht reproduzierbar bauen — Carveout.

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

**Randformen des neuen Vertrags** (`AGENTS.md` §3.12): Die Prüfsumme je Zielplattform ist ein neuer Vertrag (zwei Builds desselben Quellstands, dieselbe Prüfsumme, `SPEC-035`). Der Architect nennt und entscheidet vor dem ersten Code-Commit seine Randformen (etwa fehlende oder unbekannte Zielplattform, abweichende Toolchain, Aufrufarten von `make` nach §3.12, je Randform der Ort in der Spezifikation oder unter `harness/vertraege/`); der Implementer entscheidet keine still. Zu jeder Zusage steht eine Mutation (§3.10).

**Risiken:**

- Handbuch und Releasing-Dokument beschreiben Verhalten vor der Software; Abweichungen zeigen sich erst im Bau — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-17): Das Register ist seither gefüllt, aber kein Eintrag nennt Container, Image oder reproduzierbaren Build; einschlägig sind `BEO-REPO/verhalten-nur-unter-linux-geprueft` und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` für Aussagen des Handbuchs über Umgebungen, sowie `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` für die Schätzung in §4.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (Dockerfile mit Fragmenten, Dokumentation).

**Nachgezählt** mit der Sendung aus `slice-harness-meldungskatalog-gate` (2026-10-17, `AGENTS.md` §3.13; Gate und Handbuch-Tabellen in §1, Marke und Schätzung in §4): zwei Liefer-Punkte (§2: Image mit Smoke und Handbuch, reproduzierbarer Build mit Betriebsdokumentation), unverändert; Schichten: zwei (Dockerfile mit Fragmenten, Dokumentation), die Prüfung der Reproduzierbarkeit und der Smoke gehören zur ersten; die Bindung ist kein Liefer-Punkt, das Gate hält sie. Die Vorgabe der Marke steht beim Architect.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
