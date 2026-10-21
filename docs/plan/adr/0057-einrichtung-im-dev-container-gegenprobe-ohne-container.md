# ADR-0057: Die Einrichtung im Dev-Container ist ein Werkzeug im Repo; ihre Gegenprobe fährt ohne Container, der Lauf gegen einen echten Container ist ein Nachweis

**Status:** Accepted

**Datum:** 2026-10-21

**Autor:** Architect; angenommen durch Entscheidung des Nutzers am 2026-10-21 (wörtlich, vom Koordinator übermittelt): „ok deine Empfehlungen“ auf eine Vorlage mit drei Punkten, deren erster die Annahme dieser ADR war. Zwei Hälften waren davor schon entschieden: `ssh` als Voraussetzung am Host für dieses Ziel ist Entscheidung des Nutzers vom 2026-10-17 (auf die Frage, ob `ssh` als Voraussetzung in Ordnung ist, antwortete er „ja“), und die Obergrenze von zehn Minuten für jede Gegenprobe ist seine Entscheidung vom 2026-10-21 ([ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md)).

**Bezug:** [ADR-0045](0045-container-start-als-werkzeug-im-repo.md) (Ort eines solchen Ziels im Repo), [ADR-0047](0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (Aufteilung der Gegenproben in Container und Host), [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md), [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) (Obergrenze 600 s), [ADR-0051](0051-veroeffentlichung-auf-github-actions.md) (Nachweis als Handlung, [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 20), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes). Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`.

**Schärft:** [`SPEC-069`](../../../harness/vertraege/SPEC-069-container-einrichten.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`make container-start` liefert einen laufenden Dev-Container ([ADR-0045](0045-container-start-als-werkzeug-im-repo.md), [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md)). Was darin liegt, legt niemand an: Bare-Repo, eingelieferter `main`, Arbeitsklon, Identität des Klons, Commit-Träger und Zweig des Slice sind neun Befehle aus `docs/developer/devcontainer.md` §6 und §9, von Hand zu befolgen; das Dokument veraltet still gegen das, was funktioniert (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Die Arbeit im Container ist der Standardweg, weil ein Implementer-Lauf im Hauptbaum des Hosts den Stop-Hook über Stunden sperrt (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`).

Damit stellt sich eine Frage, die es bei `container-start` nicht gab. Seit dem 2026-10-21 hat jede Gegenprobe eine Obergrenze von **600 s** für den vollen Lauf ([ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md), Entscheidung des Nutzers). Die vergleichbare `container-start-gegenprobe` braucht mit ihrer Gruppe Docker **19 min 38 s** (gemessen und belegt bei der Lieferung ihres Teils II, Beleg in der Planung) — das Dreifache der Obergrenze. Gemessen am 2026-10-21 im Dev-Container (20 Kerne, Last 4,8):

| Größe | Wert |
|---|---|
| Fall mit Attrappen: Welt (Host-Repo mit einem Commit) | 14 ms |
| Fall mit Attrappen: Gutfall (`init --bare`, `push`, `clone`, zwei `config`, `switch`) | 43 ms |
| Fall mit Attrappen: Abbruchfall (nur Vorbedingungen) | 8 ms |
| Aufräumen je Fall | 3 ms |
| Fall gegen einen echten Container (vier bis sechs `ssh`-Runden, `make` im Container, `traeger-fetch` mit Netz) | 3 bis 25 s |
| Erststart eines Scratch-Gerüsts bei warmem Cache | 65 s (gemessen bei der Umstellung des Gerüsts auf `u-boot` 0.10.0) |

Eine Mutation fährt die Fälle ihrer Gruppe erneut. Bei rund 65 Fällen und rund 120 Mutationen ergibt das mit Attrappen 200 bis 400 s bei `GP_PARALLEL` 4 und mit einem echten Container 20 000 s und mehr. Eine Gegenprobe mit Container gehörte zudem nach [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 3 in die Liste für den Host und liefe damit weder in `make gegenproben` noch in dem Dev-Container, in dem der Implementer arbeitet.

## Entscheidung

Wir wählen **Option B: ein Ziel im Repo, eine Gegenprobe mit Attrappen ohne Container, und ein Nachweis-Lauf gegen einen echten Container als Handlung.**

1. **Ort und Art.** `make container-einrichten` liegt im Repo (Fragment unter `harness/mk/`, Skript unter `tools/harness/`), ist kein Gate und steht im Abschnitt *Werkzeuge* von `harness/README.md` — dieselbe Form wie `container-start` ([ADR-0045](0045-container-start-als-werkzeug-im-repo.md)) und `hooks-install`. Es setzt einen laufenden Container voraus und startet keinen.
2. **Host-Voraussetzungen** sind `ssh` (Entscheidung des Nutzers vom 2026-10-17) und `git`; `u-boot`, die `devcontainer`-CLI und `node` gehören `container-start` und werden von diesem Ziel nicht geprüft und nicht gebraucht.
3. **Die Gegenprobe fährt ohne Container.** `container-einrichten-gegenprobe` steht in `GEGENPROBEN_CONTAINER`, nutzt Attrappen von `ssh` (führt den Befehl lokal aus, protokolliert Argumente und Standardeingang) und von `make`, dazu **echtes `git`** auf lokalen Pfaden in Welten mit einem Commit. Geschätzte Zeit des vollen Laufs: 200 bis 400 s, Obergrenze 600 s. Sie ist **ab ihrem Anlegen aufgenommen** im Sinn von `SPEC-060` Punkt 16: Sie trägt die Fälle von dessen Punkt 16 (9), und der Slice trägt sie in der Tabelle dort als `aufgenommen` ein.
4. **Der Lauf gegen einen echten Container ist ein Nachweis, kein Sensor** — dieselbe Lesart wie `make homebrew-nachweis` ([`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 20): kein Gate, keine Gegenprobe, in keiner der beiden Listen, kein Fall und keine Abdeckungs-Deklaration; es urteilt der Mensch, der die Ausgabe liest. Er hat **keinen eigenen Treiber**: Jeder Schritt ist ein `make`-Ziel oder ein Kommando, und ein Skript darüber wäre ein Artefakt mit eigenen Zusagen, ohne einen Schritt hinzuzufügen. Seine Schritte, sein Präfix `pgr-ce-`, sein Aufräumen und die Form seines Belegs stehen in `SPEC-069` Punkt 18; der Beleg steht in §7 des Slice-Plans, der ihn führt.
5. **Was das Ziel nicht anfasst.** Es setzt die Identität nur `--local` im Klon (nie `--global`, nie `--system`, nie ein Name aus dem Skript), schreibt am Host allein das Remote, nie `core.sshCommand` (es nutzt `GIT_SSH_COMMAND` für seinen eigenen Push), nie `known_hosts`, nie `~/.ssh`, und ruft `docker`, `u-boot` und `devcontainer` nie auf. Die Entscheidung über den Hostschlüssel bleibt beim Nutzer (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Jeder `ssh`-Aufruf des Ziels trägt `BatchMode=yes`, setzt aber keine Option für den Hostschlüssel; wer eine eigene `known_hosts`-Datei braucht, gibt sie in `PGR_SSH_OPTS`.

Die Regeln im Einzelnen (Eingaben, Zustand, Schritte, Randformen, Reihenfolge der Gründe, Exit-Codes, Aufrufarten, Grenzen) stehen in [`SPEC-069`](../../../harness/vertraege/SPEC-069-container-einrichten.md) (`AGENTS.md` §3.8).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Gegenprobe mit Gruppe Docker wie `container-start-gegenprobe` | ein Artefakt; der ganze Weg an echten Werkzeugen geprüft | 20 000 s und mehr gegen eine Obergrenze von 600 s ([ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md)); Liste für den Host, also kein Lauf im Dev-Container und nicht in `make gegenproben` |
| **B — Gegenprobe mit Attrappen, Nachweis als Handlung** | 200 bis 400 s; läuft im Dev-Container und in `make gegenproben`; echtes `git`, also echte Fehlerformen; der Container-Lauf bleibt als Beleg erhalten | der Nachweis ist nicht erzwungen, er hängt am Bericht des Slice (wie [ADR-0047](0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) es für `container-start-gegenprobe` schon hinnimmt); zwei Artefakt-Arten statt einer |
| C — Gegenprobe mit einem einzigen Container-Fall je vollem Lauf (nicht je Mutation) | ein Artefakt, der Weg einmal echt geprüft | 100 bis 150 s allein für Start und Abbau (20 bis 25 % des Budgets) für **einen** Fall ohne Mutation; die Gegenprobe wird zum Host-Ziel und verliert die Fahrt im Dev-Container; ein Fall ohne Mutation ist nach `AGENTS.md` §3.10 kein Beleg |
| D — Eintrag in der Tabelle von `SPEC-060` Punkt 16 als „nicht gemessen“ mit einem anderen Plan als Nehmer | nichts zu entscheiden, die Messung kommt später | verschiebt die Frage, statt sie zu beantworten: der Nehmer müsste eine Gegenprobe beschleunigen, die mit Container nie unter 600 s kommt; dazu drei Eintragungen in einem fremden Plan (`AGENTS.md` §3.13) und das Risiko, ihn über seine Marke zu heben |
| E — Lockerung der Obergrenze für diese eine Gegenprobe | der Container-Lauf bleibt in der Gegenprobe | eine Lockerung nach `AGENTS.md` §3.6 und [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) Punkt 5, die der Nutzer entscheidet; sie kostet die Zusage, die er am 2026-10-21 gerade gegeben hat |
| F — kein Werkzeug, Dokument allein (Bestand) | nichts zu pflegen | neun Befehle je Sitzung, das Dokument veraltet still (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`), und die Arbeit im Container bleibt teurer als im Hauptbaum |

## Konsequenzen

- Positiv: Ein Aufruf richtet eine Sitzung ein; die Randformen (zweiter Aufruf, vorhandenes Bare-Repo, vorhandener Klon, abgewiesener Push, fehlendes Netz, Abbruch) sind vor dem Code entschieden und an echtem `git`-Verhalten geprüft. Die Gegenprobe läuft in dem Container, in dem gearbeitet wird, und in `make gegenproben`.
- Negativ: Der Nachweis gegen einen echten Container ist nicht erzwungen; er hängt am Bericht des Slice, und das Review prüft den Beleg. Zwei Artefakt-Arten (Gegenprobe und Nachweis) statt einer.
- Akzeptiertes Negativ: Die Gegenprobe belegt nicht, dass ein echtes `ssh` die Argumente annimmt und dass `make hooks-install` und `make traeger-fetch` im Klon gelingen (`SPEC-069` Punkt 17). Grund: Beides braucht einen Container und liegt damit über der Obergrenze; der Nachweis belegt es einmal je Lieferung.
- Akzeptiertes Negativ: Bricht der Nachweis ab, bleibt ein Container des Präfixes `pgr-ce-` stehen (`SPEC-069` Punkt 18). Grund: Ohne Treiber fängt ihn keine Falle; das Präfix macht ihn auffindbar, und das Aufräumen steht mit seiner Ausgabe im Beleg.
- Akzeptiertes Negativ: Die Zahlen oben gelten für Linux und für die gemessene Last. Grund: Die Gegenproben laufen im Gegenproben-Container (`SPEC-060` Kopf, Plattform), und die Obergrenze gilt auf ruhigem Host (`SPEC-060` Punkt 16 (1)).
- Folgepflicht: `harness/README.md` nennt Ziel und Gegenprobe erst mit der Lieferung (`AGENTS.md` §4); die Tabelle in `SPEC-060` Punkt 16 nimmt die Zeile der neuen Gegenprobe mit ihrem Anlegen auf; `docs/developer/devcontainer.md` §6 und §9 nennen das Ziel statt der Einzelbefehle. Lastenheft unverändert: Das Ziel ist ein Werkzeug der Prüfumgebung ([`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), Messmethode 4), kein Anspruch des Produkts.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, echtes `git`, Attrappen von `ssh` und `make` | Eingaben, Werkzeuge, Zustand, die sieben handelnden Schritte, Reihenfolge der Gründe, Abbruch, Aufrufarten; je Zusage und je Element eine Mutation | `make container-einrichten-gegenprobe` |
| bash, make, Attrappen | Obergrenze 600 s, Zeile mit Zeit und Last, `GP_BUDGET`, Arbeitsordner entfernt (`SPEC-060` Punkt 16 (9)) | dasselbe Ziel |
| Mensch | Lauf gegen einen echten Container, Beleg in §7 des Slice-Plans (`SPEC-069` Punkt 18) | kein Target (Handlung) |

## Re-Evaluierungs-Trigger

Wenn der Nutzer die Obergrenze aus [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) ändert; wenn ein Scratch-Container in Sekunden statt in einer Minute startet (dann ist Option C neu zu prüfen); wenn der Nachweis zweimal in Folge nicht gefahren wird (dann trägt der Beleg nichts, und die Zusagen aus `SPEC-069` Punkt 17 *akzeptiertes Negativ* sind neu zu entscheiden); wenn `u-boot` und die `devcontainer`-CLI in jedem Dev-Container liegen.
