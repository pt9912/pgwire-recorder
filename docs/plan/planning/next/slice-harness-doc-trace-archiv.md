# Slice slice-harness-doc-trace-archiv: Archivierte Slices tragen ihre Anforderungen in der Matrix (`make archiv-bezug`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; eingesammelt wird er von der nächsten Welle-Closure.
Angelegt am 2026-10-19 nach Befund des Nutzers: `make doc-trace` zeigt neun Waisen, seit `welle-v1-abschluss` archiviert ist.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Prüfbarkeit durch Werkzeuge des Harness, Messmethode 4; das Werkzeug dieses Slice ist keine Zusage des Produkts). Betroffen in der Matrix von `make doc-trace`, nicht geändert: `LH-FA-04`, `LH-FA-14`, `LH-FA-15`, `LH-FA-19`, `LH-QA-04`, `LH-RB-02` (Waisen laut §1; geliefert wird keine davon). Keine aktive ADR; ob die Form des Stubs eine ADR verlangt, entscheidet der Architect (§6).

**Berührte Spec-Stellen:** `SPEC-060` Punkte 2 und 3 (die Liste, an der die Gegenprobe hängt) · `SPEC-068` (neu; nächste freie Nummer nach `SPEC-067`, `SPEC-065` und `SPEC-066` sind gehalten; die Nummer vergibt der Architect; Vertrag in `harness/vertraege/`, Ort nach `MR-002`, vor dem ersten Code-Commit)

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Befund (gemessen am 2026-10-19, `make doc-trace` im Klon, d-check-Image aus `DCHECK_REF`, Quelltext von d-check gelesen):**
`trace.slices.dir` ist `docs/plan/planning/done`, und d-check liest dieses Verzeichnis **rekursiv** (`DiscoverFiles`, `walkMarkdown`):
die 97 Stubs unter `done/welle-*/` werden gefunden und tragen ihren Namen als Slice. Sie zählen aber nur für die Anforderungs-Kennungen,
die im Text der Datei stehen, und ein Stub trägt keine (nur Archiv-Zeiger, `BEO-`, `ADR-` und Folge-Slices in `Hervorgegangen:`). Eine Probe
mit `**Bezug:** LH-FA-14, LH-FA-15` in einem Stub senkte die Waisen von neun auf sieben (Änderung zurückgenommen). Die Konfiguration
`trace.slices` ist also nicht die Ursache, und die Zeile in `harness/README.md` §Werkzeuge („ein archivierter Slice unter
`done/<welle-id>/` zählt nicht“) beschreibt sie falsch. Die neun Waisen: `LH-FA-04`, `-08`, `-14`, `-15`, `-16`, `-19`,
`LH-QA-04`, `-05`, `LH-RB-02`. Aus den Köpfen der Volltexte in den vier `archiv.zip` (Schätzung des Planners mit `unzip` und `awk`, nicht mit dem
Werkzeug gemessen) bleibt danach eine Waise: `LH-FA-19` (Homebrew; `slice-v1-abschluss-homebrew` steht in `next/`, die Waise ist richtig).

**Ziel:** `make archiv-bezug` liest aus `docs/plan/planning/done/<welle-id>/archiv.zip` den Kopf des Volltexts jedes Slice, sammelt die
Kennungen `LH-FA-<n>`, `LH-QA-<n>`, `LH-RB-<n>` aus der Zeile `**Bezug:**` und schreibt sie als Zeile `**Bezug:**` in den Stub
(idempotent, sonst ändert es nichts); angewandt auf die vier vorhandenen Archive zählt `make doc-trace` die archivierten Slices wieder, ohne
dass `.d-check.yml` sich ändert.

**Aufnahme:** Die Liefer-Punkte sind (1) das Werkzeug mit seinem Vertrag, (2) seine Gegenprobe, (3) die Anwendung auf die 97 Stubs mit der
Berichtigung der Zeile in `harness/README.md`. Die Konfiguration `trace.slices` bleibt unverändert (Befund oben: Rekursion ist da).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung der Archivierung in `ai-harness-init` (`archive-welle` soll die Zeile `Bezug:` selbst schreiben, damit künftige Stubs sie
  tragen) — anderer Vorgang (Klasse 3): Das Werkzeug gehört einem anderen Repo. Der Slice stellt dort einen **Antrag** (Text im Plan,
  Weitergabe durch den Nutzer): `archive-welle` übernimmt die Kennungen der Anforderungen (Präfix aus `harness/conventions.md`) aus dem
  Kopf des Volltexts in den Stub, in `Hervorgegangen:` oder in einer eigenen Zeile. Bis dahin gilt für jede künftige Archivierung ein Lauf
  von `make archiv-bezug`; das ist eine Handlung, kein Sensor (siehe nächster Punkt).
- Ein Sensor, der einen Stub ohne `Bezug:` rot färbt, oder die Aufnahme von `archiv-bezug` in `GATE_CHECKS` — anderer Vorgang: Das Werkzeug
  schreibt, es prüft nicht; ein Sensor wäre ein vierter Liefer-Punkt und eine Zusage über die Archivierung eines fremden Werkzeugs.
  Ob der Architect einen Modus „prüfen ohne zu schreiben“ im Werkzeug für nötig hält, ist Frage 1 in §6; ohne Antwort bleibt es beim
  Schreiben. Eine Adresse gibt es dafür nicht: Kein Slice in diesem Repo nimmt es an.
- Ein Antrag an d-check — Klasse 3: Das Werkzeug ist Bestand eines anderen Repos und liest bereits rekursiv. Als **Antrag** (Text im Plan)
  bleibt, dass `spec/spezifikation.md` von d-check die Rekursion von `trace.slices.dir` ausdrücklich zusagt; eine Änderung am Quelltext verlangt
  dieser Slice nicht.
- Die Stubs der Wellen (`welle-*-results.md` flach in `done/`) und die Review-Reports — Bestand bleibt bewusst stehen: Sie sind keine Slices,
  `file-pattern` trifft sie nicht, und die Ergebnisnotiz bleibt vollständig und flach (Modul 6, Schritt 4).
- Ein Rückschreiben von `Berührte Spec-Stellen`, ADR-Kennungen oder Bezug auf `SPEC-`/`ARC-` — Bestand bleibt bewusst stehen: Die Matrix
  braucht nur die Anforderungs-Kennungen, und jede weitere Kennung im Stub vergrößerte das Risiko, dass das Doku-Gate (Linkpflicht der
  ADR-Kennungen) gegen Stubs meldet.
- Benutzerhandbuch, `README.md`, Produkt-Code, `spec/` — Schicht-Abgrenzung: Der Slice ändert Werkzeug, Vertrag, eine Zeile der Gate-Tabelle
  und Planung; das Produkt zeigt kein neues Verhalten (`AGENTS.md` §3.11).
- Ein Vergleich, ob die Anforderung im Volltext *geliefert* oder nur *genannt* ist („beschriebenes Verhalten, nicht geändert“) — Bestand
  bleibt bewusst stehen: `d-check` zählt in lebenden Slices genauso jede Nennung im Text; das Werkzeug übernimmt die Nennung der Zeile
  `Bezug:` und urteilt nicht darüber.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas
mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Werkzeug und Vertrag):** `tools/harness/archiv-bezug.sh` mit dem Fragment `harness/mk/archiv-bezug.mk`
      (Ziel `archiv-bezug`, `.PHONY`, Rezept ohne `$(MAKE)`; Werkzeug, kein Gate, nicht in `GATE_CHECKS`) und dem Vertrag `SPEC-068`
      in `harness/vertraege/` (Eintrag im Index und in der Historie von `harness/vertraege/README.md`), der die Randformen aus §6 entscheidet.
      Aufruf: `make archiv-bezug` für alle Archive unter `docs/plan/planning/done/*/archiv.zip`, `make archiv-bezug WELLE=<welle-id>` für eines.
      Das Werkzeug schreibt nur die Stubs und committet nichts.
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `make archiv-bezug-gegenprobe` (Fragment `harness/mk/archiv-bezug-gegenprobe.mk`, Skript
      `tools/harness/archiv-bezug-gegenprobe.sh`, Liste `GEGENPROBEN_CONTAINER` nach `SPEC-060` Punkt 2, nicht `GATE_CHECKS`) in Temp-Repos mit selbst
      gebauten Archiven, ohne Docker: je Zusage und je Element einer Aufzählung von `SPEC-068` ein Fall und eine Mutation, die genau die benannten
      Fälle rot färbt (`AGENTS.md` §3.10); die Zeile im Gate-Index von `harness/README.md` §Sensors nennt Ziel, Vertrag und Bindung.
- [ ] **Liefer-Punkt 3 (Anwendung und Berichtigung):** `make archiv-bezug` ist über die vier Archive (`welle-walking-skeleton`,
      `welle-extended-query`, `welle-replay-semantik`, `welle-v1-abschluss`; 97 Stubs) gelaufen, das Ergebnis in **einem** Inhalts-Commit je
      Welle; `make doc-trace` zeigt danach die Waisen, die §7 mit Zahl und Namen nennt (erwartet: eine, `LH-FA-19`); die Zeile von `make doc-trace`
      in `harness/README.md` §Werkzeuge nennt die Regel an den Schreiber (ein Stub trägt seine Anforderungen in `**Bezug:**`; Quelle der Matrix ist
      `trace:` in `.d-check.yml`) statt der falschen Beschreibung des Verhaltens (`AGENTS.md` §3.11).
- [ ] Mutation je Zusage der Gegenprobe selbst (`AGENTS.md` §3.10), vom Implementer rot gesehen; der Bericht nennt je Zusage: Zusage · Mutation · roter Test.
- [ ] Wirkung auf bestehende Gates gemessen und in §7 belegt: `make docs-check` (Stubs mit Anforderungs-Kennungen ohne Befund, einschließlich
      der Matrix-Regeln für die Klasse Slice), `make kopf-check` (liest nur `open/`, `next/`, `in-progress/`), `make abdeckung-check` (liest `docs/user/`,
      unberührt), und der Lauf des Stop-Hooks über Positivliste und Dateiliste.
- [ ] `make gates` grün; `make archiv-bezug-gegenprobe` und `make gates-schnell-gegenprobe` (Namensregel der Gegenproben) grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-068-archiv-bezug.md` | neu | Vertrag mit den Randformen aus §6, vor dem ersten Code-Commit; Index und Historie in `harness/vertraege/README.md` |
| `tools/harness/archiv-bezug.sh` | neu | bash, netzlos; liest `archiv.zip` mit `unzip -p`, schreibt die Zeile `**Bezug:**` in jeden Stub; Pfad der Welle und Prüfung vor dem ersten Schreiben |
| `harness/mk/archiv-bezug.mk` | neu | Ziel `archiv-bezug`, `.PHONY`, kein Anschluss an `GATE_CHECKS` |
| `tools/harness/archiv-bezug-gegenprobe.sh` | neu | Fälle und Mutationen je Zusage in Temp-Repos; nach `SPEC-068` — Happy/Boundary/Negative |
| `harness/mk/archiv-bezug-gegenprobe.mk` | neu | Ziel `archiv-bezug-gegenprobe`, `GEGENPROBEN_CONTAINER += …` |
| `harness/README.md` | update | Zeile für `make archiv-bezug` unter Werkzeuge, Zeile für `make archiv-bezug-gegenprobe` unter Sensors, Berichtigung der Zeile `make doc-trace` |
| `docs/plan/planning/done/welle-*/slice-*.md` (97 Stubs) | update | je Stub eine Zeile `**Bezug:**` hinter `**Hervorgegangen:**`, erzeugt vom Werkzeug, in einem Commit je Welle |

Tests und Gates: neue Gegenprobe `archiv-bezug-gegenprobe`. Reihenfolge: Vertrag (Architect), Werkzeug, Gegenprobe gegen es, Gate-Verdrahtung
und README-Zeilen, zuletzt die Anwendung auf die Stubs (damit der Beleg in §7 den Lauf des fertigen Werkzeugs zeigt).

**Optional — Ansatz als Liste:**

- Der Kopf eines Volltexts endet an der ersten Zeile `---`; die Zeile `**Bezug:**` kann über mehrere Zeilen laufen (bis zur nächsten Zeile, die
  mit `**` beginnt, oder bis zum Kopfende). Aus diesem Bereich liest das Werkzeug Tokens der Form `LH-(FA|QA|RB)-[0-9]+` — rein lexikalisch,
  wie d-check selbst liest (`LH-(?:FA|QA|RB)-\d+`); ein Suffix wie `.e` in `LH-FA-05.e` gehört nicht zur Kennung.
- Der Stub trägt Kennungen, nicht Adressen (Zitier-Form des Stub-Vorlage): die Zeile enthält die Kennungen **ohne** Link.
- Antrag an `ai-harness-init` und an d-check stehen als Text in §1 (Klasse 3), nicht als Arbeit.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): wellenlos und von den Release-Slices unabhängig; er blockiert kein Gate und keinen Schritt. Die Stelle in der
Reihenfolge ist Entscheidung des Nutzers. WIP-Limit 1 je Lauf: Läuft ein zweiter Implementer parallel zu `slice-erster-release-veroeffentlichung`,
trägt `Verantwortlich:` Person und Zweig (`modul-05`, parallel). Die Basis der Metrik ist der Stand des Hauptzweigs vor dem Anspruch.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von
  `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness`, Tests und Gegenprobe eingeschlossen, **über 800**; die 97 Stubs
  unter `docs/plan/planning/` zählen nicht in die Summe (Daten, vom Werkzeug erzeugt, je Stub eine Zeile; ihre Zahl prüft §2 Punkt 3). Basis: fd958f07 (letzter Plan-Commit);
  beim Anspruch bestätigt der Implementer sie oder ersetzt sie
  durch den Stand des Hauptzweigs vor dem Anspruch. **Meldeschwelle 650:** Erreicht die Summe 650, hält der Implementer an und meldet dem Planner,
  bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 480 bis 640, Mitte 560, 70 % der Marke.
  Aufbau: Skript rund 110, Fragmente rund 30, Vertrag rund 110 bis 140 (Zeilen, `harness/` zählt mit), README-Zeilen rund 6, Gegenprobe 230 bis 340
  (das Doppelte des Produktcodes, `AGENTS.md` §3.14; Fälle: Zip fehlt, beschädigt, Stub ohne Eintrag, leerer Kopf, mehrere Anforderungen, Suffix,
  mehrzeiliges `Bezug:`, vorhandene Zeile gleich und verschieden, `WELLE=`, Aufrufarten von `make`). Faktor Tests zu Produktcode: 2 bis 3
  (Beleg: `slice-harness-ci-workflows`, Roadmap-Drift-Log 2026-10-18, Messung 899 gegen Schätzung 545 bis 695); eine eigene Messung für
  eine Gegenprobe über einem bash-Werkzeug dieser Größe gibt es nicht (keine Messung).
- `in-progress` → `next`: Der Architect findet, dass der Stub neben `Bezug:` weitere Kennungen tragen muss (ADR, `SPEC-`), oder dass die Form
  des Stubs eine Adaption nach `harness/conventions.md` verlangt (§6, Frage 3) — dann gehört die Entscheidung vor den Code, und der Slice geht
  zurück, wenn sie den Umfang über die Marke hebt.
- `in-progress` → `open`: `make docs-check` meldet gegen die Stubs mit Anforderungs-Kennungen einen Befund, den eine Änderung von `.d-check.yml`
  verlangt (Änderung des Gates, `AGENTS.md` §3.6), oder das gepinnte d-check liest nicht mehr rekursiv — der Slice wartet auf die Entscheidung.

**Basis der Metrik:** fd958f07 (letzter Plan-Commit, Move nach `next/`); beim Anspruch bestätigt der Implementer sie oder ersetzt sie durch den Stand des Hauptzweigs vor dem Anspruch.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make doc-trace` zeigt nach dem Lauf die in §7 genannte Zahl von Waisen (erwartet: eine, `LH-FA-19`); Review-Report liegt vor;
`make gates` grün; Closure-Notiz mit Lerneintrag (Form: neuer Sensor `make archiv-bezug-gegenprobe`, oder benannte Spec-Lücke zu `archive-welle`:
Der Stub verliert die Anforderungs-Kennung, die die Matrix braucht).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Vom Architect entschieden am 2026-10-19 (Ort: `SPEC-068`):** (1) Modus „prüfen“: **nein** (`SPEC-068` Punkt 11 a). (2) Wirkung auf die Matrix als Zusage mit Docker-Fall: **nein**; einmalige Messung in §7 (Punkt 11 b). (3) Form: eigene Zeile `**Bezug:**` hinter `**Hervorgegangen:**`, keine Adaption `MR-<NNN>` (die Kennungen der Anforderungen gehören zu denen, die den Vorgang überlebt haben, Baseline `modul-06-roadmap.md` Schritt 4), Vorlage unverändert (Kopf von `SPEC-068`). (4) Eine ADR: **nicht nötig** (kein Gate gelockert, keine Architekturregel; Gründe im Kopf von `SPEC-068`). (5) Nummer `SPEC-068`. (6) Alle Randformen unten: entschieden in `SPEC-068`, je Zeile mit Punkt. (7) `unzip` und `zip`: in Dev- und Gegenproben-Container gemessen vorhanden (`/usr/bin/unzip` 6.0, `/usr/bin/zip`); der Lauf prüft das Werkzeug vorab (Exit 2, kein Überspringen). (8) Wirkung auf `docs-check`, `kopf-check`, `abdeckung-check` und den Stop-Hook: Vorhersage in `SPEC-068` Punkt 11 c, Messung in §7.

Die Fragen unten bleiben als Herkunft stehen:

1. *Modus „prüfen“:* Soll das Werkzeug neben dem Schreiben einen Lauf ohne Änderung bieten, der mit Exit 1 endet, wenn ein Stub die Zeile nicht
   trägt (kein Gate, kein Anschluss an eine Kette)? Vorschlag des Planners: nein (ein Sensor ohne Anschluss prüft nichts, und `archive-welle`
   upstream soll die Zeile schreiben); der Architect entscheidet, mit Folge für Zählung und Marke.
2. *Wirkung auf die Matrix als Zusage:* Die Gegenprobe läuft ohne Docker. Soll ein Fall zeigen, dass d-check die Zeile liest (Image aus
   `DCHECK_REF`, wie `range-konfiguration-gegenprobe`)? Vorschlag des Planners: nein; der Satz in `harness/README.md` nennt nur die Regel an den
   Schreiber (`AGENTS.md` §3.11), die Wirkung steht als einmalige Messung in §7. Sagt der Architect ja, wächst die Gegenprobe um einen Docker-Fall
   und gehört dann zu den Gegenproben mit Image.
3. *Form des Stubs:* Die Stub-Vorlage der Baseline führt „vier Dinge und sonst nichts“ (Identität, Archiv-Zeiger, Zustand, die Kennungen, die den
   Vorgang überlebt haben). Ist eine eigene Zeile `**Bezug:**` eine Adaption nach `harness/conventions.md` (`MR-<NNN>`, mit ADR) oder zählen die
   Anforderungs-Kennungen zu „den Kennungen, die überlebt haben“ (dann reicht eine Entscheidung im Vertrag)? Alternative: die Kennungen in
   `Hervorgegangen:` schreiben. Vorschlag des Planners: eigene Zeile (kein Parsen der von `ai-harness-init` geschriebenen Zeile, einfache
   Idempotenz) und keine MR, weil die Kennungen zur Zitier-Form des Stubs gehören; der Architect entscheidet, und die Vorlage der Baseline
   (`archiv-stub-slice.template.md`) bleibt unverändert.
4. *Eine ADR?* Neue Werkzeuge dieser Art tragen eine ADR mit Entscheidung und Gründen ([ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md)); der Architect entscheidet, ob
   „Anforderungen im Stub statt Konfiguration von d-check“ eine ist (`AGENTS.md` §3.8: die ADR zeigt mit `Schärft:` auf `SPEC-068`).

**Randformen des Werkzeugs** (`AGENTS.md` §3.12; Ort der Entscheidung: `SPEC-068`; Vorschlag des Planners je Zeile, entschieden ist erst, was der
Architect in den Vertrag schreibt):

- **Entschieden:** `SPEC-068` Punkt 4 (Leerwert `**Bezug:** — keine —`, zweiter Lauf ohne Änderung). *Leerer Kopf* (kein `**Bezug:**` im Volltext, oder kein Token): Vorschlag: die Zeile `**Bezug:** — keine —` (Leerwert der Vorlage), damit der Stub
  als verarbeitet gilt und der zweite Lauf nichts ändert; Fall und Mutation.
- **Entschieden:** `SPEC-068` Punkt 3 (alle Stubs geprüft vor dem ersten Schreiben, Exit 1, alle Befunde genannt). *Stub ohne Eintrag im Archiv* (Datei `slice-*.md` im Welle-Verzeichnis, aber kein Eintrag gleichen Pfads in `archiv.zip`): Vorschlag: Exit 2 vor dem
  ersten Schreiben, die Meldung nennt den Stub; keine Datei der Lauf-Menge ändert sich.
- **Entschieden:** `SPEC-068` Punkt 2 (Exit 2, Reihenfolge: Form der Welle, Verzeichnis, `unzip`, Archiv fehlt, Archiv beschädigt). *Archiv fehlt* (Welle-Verzeichnis ohne `archiv.zip`; `WELLE=` unbekannt) und *Zip beschädigt* (`unzip -t` scheitert): Vorschlag: Exit 2 vor dem ersten
  Schreiben, ohne Änderung, mit Namen; Reihenfolge der Gründe (`AGENTS.md` §3.12): Eingabe der Welle, Archiv fehlt, Archiv beschädigt, Stub ohne Eintrag.
- **Entschieden:** `SPEC-068` Punkt 5. *Mehrere Anforderungen:* sortiert, ohne Duplikate, durch `, ` getrennt; Vorschlag der Ordnung: `LH-FA`, `LH-QA`, `LH-RB`, innerhalb nach Zahl
  (nicht bytewise: `LH-FA-9` vor `LH-FA-10`).
- **Entschieden:** `SPEC-068` Punkt 4. *Mehrzeiliges `Bezug:` und Suffixe* (`LH-FA-05.e`, `LH-QA-03 und -04`): Vorschlag: gelesen werden nur volle Tokens `LH-(FA|QA|RB)-[0-9]+`; die Kurzform `-04`
  ist keine Kennung (d-check liest sie ebenso nicht); die Kennung vor dem Suffix zählt.
- **Entschieden:** `SPEC-068` Punkt 6. *Schon vorhandene Zeile:* gleich → keine Änderung (Bytes, Zeitstempel); abweichend → Vorschlag: nicht überschreiben, Exit 1 mit Stub und beiden Werten
  (eine von Hand gepflegte Zeile geht nicht verloren); zwei Zeilen `**Bezug:**` im Stub → Exit 1.
- **Entschieden:** `SPEC-068` Punkte 7 und 8 (CRLF im Stub ist Befund; Volltext mit CRLF nicht zugesagt). *Stelle der Zeile:* hinter `**Hervorgegangen:**`, vor einer Leerzeile; fehlt `**Hervorgegangen:**`: Vorschlag Exit 1 (Stub hat nicht die Form der Vorlage).
  Letzte Zeile ohne Zeilenumbruch: Zeile wird angehängt; CRLF im Stub: nicht zugesagt (`BEO-REPO/zeilenende-crlf-im-gelesenen-wert`), Entscheidung des Architect.
- **Entschieden:** `SPEC-068` Punkt 9 (ein Aufruf nach dem anderen; gleichzeitige nicht zugesagt). *Aufruf-Voraussetzung:* ein Aufruf nach dem anderen im Arbeitsbaum; der Lauf braucht keinen sauberen Baum und committet nicht. Gleichzeitige Aufrufe sind nicht
  zugesagt (Entscheidung offen).
- **Entschieden:** `SPEC-068` Punkt 10 (Fälle `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`). *Aufrufarten des Treibers* (`AGENTS.md` §3.12): das Werkzeug liefert kein Ergebnis, das ein Hook oder ein Gate glaubt; Vorschlag: `make -n` druckt und schreibt
  nichts (Rezept ohne `$(MAKE)`), `-q` und `-t` ändern nichts, `-i` ändert am Exit nichts; der Architect bestätigt oder streicht die Nennung mit Grund.
- **Entschieden:** `SPEC-068` Kopf, Plattform. *Plattform* (`AGENTS.md` §3.12, `LH-QA-03`): geprüft nur unter Linux mit bash und `unzip`; macOS und Windows als Grenze im Vertrag (akzeptiertes Negativ). `unzip`
  ist im Dev-Container vorhanden (`/usr/bin/unzip`); ob die Container der Gegenproben es führen, belegt der Implementer vor dem Code (nicht gemessen).

**Risiken:**

- Künftige Archivierungen schreiben den Stub weiter ohne Anforderungs-Kennung, und die Waisen kehren zurück, ohne dass ein Gate es meldet (`make doc-trace` ist
  kein Gate). — **Ausgang:** — (bei Closure; Vorschlag: *weiter offen* → neues `BEO-REPO/archiv-stub-verliert-anforderungs-kennung` im Register, bis `ai-harness-init`
  den Antrag aus §1 aufnimmt)
- Die Zeile `**Bezug:**` im Stub verletzt eine Prüfung, die ich nicht gemessen habe (Matrix-Regeln des Doku-Gates für die Klasse Slice, Linkpflicht, Form-Prüfung eines
  Stubs). — **Ausgang:** — (bei Closure; Messung in §7, Lauf von `make docs-check` vor dem Commit der ersten Welle)
- Die Nennung im Kopf ist keine Lieferung (`LH-…` nur zitiert), und die Matrix zählt mehr als geliefert. — **Ausgang:** — (bei Closure; Vorschlag: *entfallen*: dieselbe
  Lesart wie bei lebenden Slices, siehe §1 letzter Ausschluss)
- Das Werkzeug ändert 97 Dateien in einem Lauf; ein Fehler in der Mitte hinterließe einen halben Stand. — **Ausgang:** — (bei Closure; Gegenmittel im Vertrag: Prüfung aller
  Stubs und Archive vor dem ersten Schreiben, Fall und Mutation)

## 7. Closure-Notiz

Wird bei der Closure gefüllt (Planner, `/close-welle`-Anweisung für wellenlose Slices). Hier stehen dann die Messungen aus §2 (Waisen vor und nach dem Lauf, Gates).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration); das Repo deklariert nur diese eine,
eine Ausdifferenzierung ist nicht Gegenstand dieses Slice.

**Vorgelagert — offene Beobachtungen sichten:** Register (`docs/plan/planning/observations/BEO-REPO/`) durchgegangen, 2026-10-19. Treffer, die den Slice berühren:
`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10: Mutation je Zusage, hier in §2 verlangt), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§3.12: Randformen
in §6), `BEO-REPO/zeilenende-crlf-im-gelesenen-wert` (CRLF im Stub, §6), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand`
(§3.14, Marke mit Metrik in §4), `BEO-REPO/meldeschwelle-ohne-artefakt` (Meldeschwelle in §4). Zähler-Stände habe ich nicht neu gezählt; der Zähler folgt aus den Dateien
im Register. Kein Eintrag nennt die Archivierung, die Matrix oder `archive-welle`: dafür keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Nachzählung (`AGENTS.md` §3.13):** Dieser Slice trägt keine Sendung in einen Nehmer ein; es gibt keine Folge-Slice-Adresse (die Anträge an `ai-harness-init` und d-check sind
Text). Liefer-Punkte: drei (Werkzeug mit Vertrag · Gegenprobe · Anwendung mit Berichtigung). Schichten: eine — Werkzeug des Harness (`tools/harness/`, `harness/`); die
Stubs unter `docs/plan/planning/` sind Planung und zählen nicht, `harness/README.md` ist Harness-Datei und keine Nutzer- oder Wartungs-Doku (`docs/user/`, `docs/maintainer/`,
`README.md` bleiben unberührt). Damit höchstens zwei Schichten und höchstens drei Punkte.

**Schätzung des Architect nach dem Vertrag (2026-10-19; Metrik von §4; Faktor der Gegenprobe wie `slice-harness-ci-workflows` gemessen 899 gegen 545 bis 695):** Skript 130 bis 170 (zwei Phasen, Kopf lesen, Zeile schreiben), Fragmente 30, README-Zeilen 6, Gegenprobe 300 bis 430 (rund 55 Fälle, 40 Mutationen, Archive mit `zip`); roh zusammen 470 bis 640; `SPEC-068` (rund 60 Zeilen Vertragstext, `harness/` zählt mit, wenn die Basis vor ihrem Commit liegt) kommt mit 60 bis 70 hinzu; zusammen **470 bis 710**, Mitte um 590; die Meldeschwelle 650 kann erreicht werden, die Marke 800 liegt darüber. Schnitt, falls sie eintritt: Fälle zu Punkt 3 (halber Stand) und Punkt 10 (Aufrufarten) in einen Folge-Slice nur dann, wenn der Planner es will; der Architect sieht keinen Schnitt vor.
