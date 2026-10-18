# Slice slice-harness-devcontainer-einrichten-zugang: `container-einrichten` löst den Zugang von Claude Code auf und reicht ihn an den Container (Quellenauflösung, `devcontainer exec`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-18 durch den Schnitt von `slice-harness-devcontainer-einrichten` (Entscheidung des Nutzers vom 2026-10-18: Der Zugang von Claude Code im Container wird nicht gebraucht; kein Login, die Hostsitzung treibt die Arbeit per `ssh` und Befehl im Container). Der Slice bleibt als Plan erhalten, wird aber nicht priorisiert (§4).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (dieselbe Frage wie im Vorgänger; die `devcontainer`-CLI kommt dort an); kein Anspruch des Produkts.
Vorgänger in dieser Reihenfolge: `slice-harness-devcontainer-einrichten-geruest` (liefert das Ziel `container-start` und schreibt `mk-override.sh`), `slice-harness-devcontainer-einrichten` (liefert das Ziel `container-einrichten`; dieser Slice erweitert dessen Skript) und `slice-harness-devcontainer-u-boot-0-10` (`done/`; Geber der Sendung „Quellenauflösung“ und der Befunde zu `remoteEnv`). Keine aktive ADR; der Architect legt eine an oder ergänzt die des Vorgängers (§6).

**Berührte Spec-Stellen:** Ergänzung des Vertrags des Vorgängers in `harness/vertraege/` (Nummer vergibt der Architect im Vorgänger); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make container-einrichten` löst den Zugang von Claude Code am Host in fester Reihenfolge auf (`ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, Zugangstoken aus der Anmeldung von `claude` am Host), meldet, welche Quelle genommen wurde (nie den Wert), bricht ohne Quelle ab, bevor es im Container etwas ändert, und reicht den Wert so an den Container, dass ein Prozess dort ihn in `claude` sieht (Randform 15). Das Dokument `docs/developer/devcontainer.md` §7 stimmt mit dem Ziel überein.

**Sendung des Vorgängers `slice-harness-devcontainer-einrichten` (eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13).** Dieser Slice nimmt an, was jener ausschließt (dessen §1, Klasse 1): die Quellenauflösung samt Randform 15 und den Entscheidungen des Nutzers vom 2026-10-17 dazu (§6), die Übergabe des Werts (`remoteEnv` wirkt nur bei `devcontainer exec`), die Meldung der Quelle ohne Wert, die Gegenprobe der Quellen mit Platzhaltern und den Abschnitt §7 des Dokuments. Die Bindung des Vorgängers ruht auf diesem Slice: Dessen Skript liest und löst keinen Zugang und hält eine Stelle vor seinem Schritt 1 frei (dessen Randform 25); `mk-override.sh` schreibt der Geber `slice-harness-devcontainer-einrichten-geruest` (dessen Bindung, eingetragen vom Planner am 2026-10-18) mit der Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}`; über deren Fortbestand entscheidet dieser Slice und ändert dafür `tools/harness/container-start.sh` und dessen Gegenprobe (Sendung des Gebers, in die Zählung in §8 eingerechnet).

**Sendung „Quellenauflösung“ von `slice-harness-devcontainer-u-boot-0-10`** (eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag) geht über den Vorgänger an diesen Slice: Der Geber (`done/`) nennt weiter den Vorgänger als Nehmer (ein Hop länger, kein Zeiger im Geber wird umgeschrieben). Übernommene Befunde des Gebers (Dokument §7 und §8, Platzhalter statt Zugangsdaten): `remoteEnv` liefert `ANTHROPIC_API_KEY` und `CLAUDE_CODE_OAUTH_TOKEN` **nur an Prozesse, die die `devcontainer`-CLI mit `devcontainer exec` startet**; in einer `ssh`-Sitzung (auch mit `SendEnv`) und bei `docker exec` sind beide nicht gesetzt; `devcontainer exec` liest den Wert bei jedem Aufruf aus der Shell des Hosts, ein Neustart des Containers ist dafür nicht nötig; `claude -p` meldet ohne Variable `Not logged in`, mit Platzhalter-Schlüssel `Invalid API key`, mit Platzhalter-Token `401 Invalid bearer token`. Dieser Slice übernimmt diese Ergebnisse aus dem Dokument und misst sie nicht neu. Den echten Lauf mit echtem Token macht der Nutzer selbst: Die Berechtigungsprüfung verbietet dem Agenten, Anmeldedaten zu materialisieren, die Tests des Ziels laufen mit Platzhaltern.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Gerüst und kein Start: Das liefert `slice-harness-devcontainer-einrichten-geruest` (Klasse 1, Vorgänger; dieser Slice ändert dort nur die Zeile in `mk-override.sh`).
- Kein Bare-Repo, kein Klon, keine Identität, kein Zweig: Das liefert `slice-harness-devcontainer-einrichten` (Klasse 1, Vorgänger; dieser Slice setzt dessen Ziel voraus und ändert nur dessen Zugangs-Stelle).
- Kein Login von Claude Code im Container (`/login`, Abo-Anmeldung per Browser) und keine Aussage, dass beide Wege (API-Schlüssel, Abo-Token) im Container funktionieren: Das ist der Lauf mit echtem Token des Nutzers (Risiko in §6).
- Keine Datei mit dem Zugang im Container oder im Image, kein Mount der Anmeldedatei: Entscheidung des Nutzers vom 2026-10-17, der Zugang kommt als Umgebungsvariable, nie dauerhaft abgelegt.
- Keine Änderung am Egress des Containers (`claude.ai`, `platform.claude.com`, `claude.com` stehen nicht in `egress.allow`): Das ist eine Entscheidung zu Sicherheit und Gerüst, kein Teil der Übergabe; meldet ein gültiges Token einen Netzfehler, trägt der Nutzer den Host nach (Dokument §11).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Dokumentation, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos; die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm.

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Quellenauflösung und Übergabe im Ziel** (`tools/harness/container-einrichten.sh`, an der freigehaltenen Stelle des Vorgängers): erfüllt die Ergänzung des Vertrags, die der Architect festlegt (Randform 15 und Reihenfolge der Gründe); die Quelle wird gemeldet, nie der Wert; ohne Quelle bricht das Ziel vor jeder Änderung ab.
- [ ] **Gegenprobe der Quellen** (in `container-einrichten-gegenprobe` des Vorgängers): je Quelle, je Reihenfolge, je Element der Aufzählungen (leer, nur Leerzeichen, fehlend) ein Fall und je Zusage eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); Platzhalter statt Anmeldedaten; kein Wert in Ausgabe, Log oder Prozessliste, soweit der Vertrag es zusagt und die Gegenprobe es prüft.
- [ ] **Dokument** `docs/developer/devcontainer.md` §7 nennt das Ziel als Weg für den Zugang und führt nur Aussagen, die ein Lauf belegt oder §13 nennt; Aussagen über den echten Login bleiben auf das beschränkt, was die Gegenprobe zeigt (`AGENTS.md` §3.11).

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/` ändern); nicht gleichzeitig mit Läufen der Scratch-Container der Gegenprobe.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-<NNN>-container-einrichten.md` (Vertrag des Vorgängers) | update (Architect) | Randform 15 und ihre Entscheidungen; Reihenfolge der Gründe; Übergabe |
| `docs/plan/adr/` (ADR des Vorgängers oder eine neue) | update oder neu (Architect) | Entscheidung zur Übergabe (`devcontainer exec`, Wert über `stdin` oder andere Form) und zur Host-Voraussetzung `devcontainer`-CLI, soweit sie nicht schon der Vorgänger trägt |
| `tools/harness/container-einrichten.sh` | update | Quellenauflösung, Übergabe, Meldung ohne Wert an der freigehaltenen Stelle |
| `tools/harness/container-start.sh` (Geber) | update | Zeile mit `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` in `mk-override.sh`: Fortbestand nach der Entscheidung des Architect |
| `tools/harness/container-einrichten-gegenprobe.sh` | update | Fälle und Mutationen der Quellen, Platzhalter |
| `harness/README.md` | update | Index nennt die Zusage der Gegenprobe, soweit der Vertrag sie erweitert |
| `docs/developer/devcontainer.md` | update | §7; §13 für Ungeprüftes |

- Die Quelle 3 (Anmeldedatei am Host) hängt an einem von Anthropic nicht zugesagten Format; der Architect entscheidet, ob sie in den Vertrag kommt (§6).
- Arbeit des Implementers im Container: wie im Vorgänger beschrieben (dort §3); dieser Slice ändert nur Dateien, die der Vorgänger angelegt hat, und setzt dessen Lieferung voraus.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten` liegen in `done/` **und** der Nutzer hat den Zugang wieder verlangt. Zurzeit gilt das Gegenteil (Entscheidung des Nutzers vom 2026-10-18: nicht gebraucht), der Slice bleibt in `open/` und wird nicht priorisiert; die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm. Wird der Zugang endgültig nicht gebraucht, geht dieser Slice mit §7 `Gegenstand: entfallen: <Grund>` von `open/` nach `done/` (Liefer-Punkte leer, jedes Risiko mit Ausgang); der Nehmer der Sendungen ist dann entfallen, und der Vorgänger hält die Stelle nicht mehr frei (Rückfrage an den Planner).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht. Schätzung (Planner, 2026-10-18; Faktor 2 für die Gegenprobe; Bestand: Teilsummen des Vorgängers vor dem Schnitt): Quellenauflösung 70 bis 110, Übergabe über `devcontainer exec` 20 bis 40, Gegenprobe 200 bis 300, Dokument 15 bis 25, Zeile in `mk-override.sh` samt Fall 5 bis 10: **310 bis 485**. Unter der Meldeschwelle 700 und der Marke 900.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder zwei Schichten, oder die Randformen von 15 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Quelle 3 und ihre Gegenprobe getrennt).
- `in-progress` → `open` (blockiert): Der Vorgänger hält die Stelle nicht frei, oder die `devcontainer`-CLI am Host ist nicht zugesagt (Entscheidung des Nutzers zu den Host-Voraussetzungen, Vorgänger Randform 17), oder `claude` nimmt keinen der beiden Wege an (Befund an den Planner, Eintrag im Register).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Der Lauf mit echtem Token bleibt Handlung des Nutzers und ist keine Closure-Bedingung (Risiko in §6).

## 6. Risiken und offene Punkte

Der Slice liefert eine **Ergänzung eines neuen Vertrags** (Zusagen: Reihenfolge, Meldung ohne Wert, Abbruch vor der Änderung); §3.12 gilt. Die Randform steht hier, entschieden wird sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); der Implementer entscheidet keine Frage still, er hält an und gibt sie dem Architect zurück.

**Randform 15 — Quelle des Zugangs von Claude Code (Quellenauflösung; übernommen aus `slice-harness-devcontainer-einrichten` vom 2026-10-18).**

**Entscheidung des Nutzers vom 2026-10-17 („ja“): `container-einrichten` löst die Quelle in fester Reihenfolge auf:** (1) `ANTHROPIC_API_KEY` der Host-Umgebung, (2) `CLAUDE_CODE_OAUTH_TOKEN` der Host-Umgebung, (3) das Zugangstoken (nie das Erneuerungstoken) aus der Anmeldung von `claude` am Host, zur Laufzeit gelesen. Durchreichen: `remoteEnv` wirkt nach dem Befund des Gebers nur bei `devcontainer exec`, nicht in `ssh`-Sitzungen und nicht bei `docker exec`; ob das Ziel dafür die `devcontainer`-CLI am Host braucht (Host-Voraussetzung, die der Nutzer entscheidet; der Vorgänger braucht sie für den Start ohnehin, wenn der Nutzer sie dort entscheidet) oder einen anderen Weg geht (zum Beispiel Wert über `stdin` in den Aufruf, Dokument §7), entscheidet der Architect mit dem Nutzer; nie dauerhaft abgelegt, nie der Wert in Ausgabe oder Log; das Ziel meldet am Ende, welche Quelle genommen wurde (nie den Wert) und bricht klar ab, wenn keine da ist. Der Nutzer verlangt, dass alles automatisierbar ist. Erhalten aus den früheren Entscheidungen des Nutzers vom 2026-10-17: Der Container muss den Zugang nicht besser schützen als der Host, kein Sicherheitsargument gegen das Durchreichen; das Token erneuert der Host (nach dem Befund des Gebers genügt der nächste `devcontainer exec`, ein Neustart des Containers ist dafür nicht nötig; ob ein laufender `claude`-Prozess den neuen Wert übernimmt, ist nicht gemessen); der Container schreibt den Zugang nie; beide Wege (API-Schlüssel und Abo) müssen im Container funktionieren; ein Mount der Datei entfällt, der Zugang kommt als Umgebungsvariable.

Vom Architect zu entscheiden und in der Spezifikation zu binden: Reihenfolge der Gründe bei mehreren Quellen (die erste vorhandene gilt, später folgende werden nicht geprüft?), leer oder nur Leerzeichen zählt als nicht vorhanden?, Übergabe an den Container (das Ziel startet den Container im Vorgänger; nur Prozesse von `devcontainer exec` bekommen die Variable), Rechte der Zwischendatei und Entfernen nach dem Start, Verhalten bei abgelaufenem Zugangstoken, Meldung ohne Wert, Fortbestand der Zeile mit `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` in `mk-override.sh` (Bindung des Vorgängers, dessen §1).

**Offene Fragen, ehrlich ungeklärt:** Quelle 3 hängt an einem von Anthropic nicht zugesagten Dateiformat (Ort, Feldnamen); ob `claude` ein Zugangstoken über `CLAUDE_CODE_OAUTH_TOKEN` annimmt (offiziell gedacht für das Token von `claude setup-token`), misst der Nutzer; ein Zugangstoken, das älter als rund 8 Stunden ist, erneuert der Host (die 8 Stunden sind die Grenze des Nutzers, nicht gemessen); die Tests des Ziels laufen mit Platzhaltern, den echten Lauf mit echtem Token macht der Nutzer.

**Weitere Risiken:**

- **Die Quellenauflösung hängt an Annahmen, die niemand zugesagt hat** (Randform 15: Dateiformat der Anmeldung von `claude`, Annahme eines Zugangstokens über `CLAUDE_CODE_OAUTH_TOKEN`). Gemessen ist nur, dass die Variable ankommt und `claude` auf Platzhalter antwortet (`401 Invalid bearer token`); ob ein gültiges Token angenommen wird, ist offen. Ein geändertes Format bricht Quelle 3 still oder laut; Tests mit Platzhaltern belegen die Reihenfolge, nicht die Annahme des Tokens durch `claude`. — **Ausgang (Closure):** weiter offen → Register (neu, Evidence dieser Slice), bis der Nutzer den Lauf mit echtem Token gemacht hat; die Aussage im Dokument beschränkt sich auf das, was die Gegenprobe zeigt (`AGENTS.md` §3.11).
- **Slice länger als rund 8 Stunden** (Lebensdauer des Zugangstokens). Das ist ein Erneuern des Tokens am Host (der nächste `devcontainer exec` liest es), kein Fehler des Werkzeugs. — **Ausgang (Closure):** entfallen als Risiko des Slice, weil das Dokument das Erneuern nennt; ein Fund im Betrieb geht ins Register.
- **Wert im Klartext in der Prozessliste des Hosts** (Dokument §7: `devcontainer exec` ruft `docker exec -e <Name>=<Wert>`, `ps` zeigt ihn). — **Ausgang (Closure):** entfallen, wenn der Vertrag die Form der Übergabe festlegt und die Gegenprobe zeigt, was sie zusagt; sonst weiter offen → Register.
- **Der Slice wird nicht gebraucht** (Entscheidung des Nutzers vom 2026-10-18). — **Ausgang (Closure):** entfallen, wenn der Slice mit `Gegenstand: entfallen: <Grund>` nach `done/` geht (§4); bleibt er in `open/`, hat er keinen Ausgang zu tragen.
- **Konflikt mit dem Skript des Vorgängers** (dieselbe Datei, dieselbe Gegenprobe). — **Ausgang (Closure):** entfallen, wenn der Vorgänger in `done/` liegt (Start-Trigger §4).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der die Sendung des Vorgängers einträgt):** Liefer-Punkte 3 (Quellenauflösung und Übergabe, Gegenprobe der Quellen, Dokument §7). Schichten 2 (Harness-Werkzeug mit Gegenprobe; Nutzer- und Wartungs-Doku `docs/developer/`; Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 310 bis 485 (mit der Bindung des Gebers `slice-harness-devcontainer-einrichten-geruest`, nachgezählt am 2026-10-18 im Commit, der sie einträgt); Meldeschwelle 700, Marke 900. Der Nehmer liegt nicht über der Größenregel; ein Schnitt unterbleibt. Die Sendung „Quellenauflösung“ von `slice-harness-devcontainer-u-boot-0-10` liegt über den Vorgänger hier; `grep` nach der Kennung des Vorgängers findet ihn in §1, §3, §4 und §6 dieses Slice.
