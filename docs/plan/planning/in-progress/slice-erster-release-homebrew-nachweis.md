# Slice slice-erster-release-homebrew-nachweis: Homebrew-Nachweis (Abnahmeszenario 11)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-19.a` · `SPEC-042` · `SPEC-067` (Punkt 19 (7), Nachtrag; Punkt 20, Nachweis)

**Verantwortlich:** pt9912 (Zweig impl-nachweis)
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Nach dem ersten stabilen Release liegt die Formel im Tap `pt9912/homebrew-pgwire-recorder`, und die dokumentierten Befehle installieren das Werkzeug. **Stand nach dem Lauf** (Implementer, 2026-10-21, §3.9): Belegt ist die Kombination **Linux/amd64** (Beleg in §7); macOS (`amd64` und `arm64`) und Linux auf `arm64` sind **Grenze** in `LH-FA-19.a` und im Handbuch — die Zusage gilt je Kombination einzeln (§6 Randform 9), und das Ziel behauptet seit dieser Nachführung keine Plattform, die kein Beleg trägt.

**Übernommen aus `slice-erster-release-veroeffentlichung`** (dort §1, Abgrenzung): die Homebrew-Formel im Tap und ihr Nachweis; das sind Ziel und DoD-Punkt 1 oben.

**Übernommen aus `slice-v1-abschluss-homebrew`** (dort §1 und §6, `AGENTS.md` §3.13): (1) der Abschnitt zu Homebrew in §2 *Installation* des Benutzerhandbuchs und die Zeile in `README.md`, jeweils nur mit Befehlen, die gegen den echten Tap gelaufen sind; (2) die **Formel des ersten stabilen Tags** im echten Tap: Sie entsteht durch den Job `homebrew` des Laufs jenes Tags (`SPEC-067` Punkt 19); schlägt der Job fehl, ist der Wiederanlauf des Jobs oder der Nachtrag von Hand mit den Token in der Umgebung der Weg (Punkt 19 (7), Handlung des Betreibers). Für `v0.1.0` gibt es keine Formel: Seine Binaries tragen das alte Schema (Entscheidung des Nutzers vom 2026-10-20, [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md)); (3) die Belege, die die Gegenprobe jenes Slice nicht liefert: das Schreibrecht des Tap-Tokens, die Annahme der Formel durch Homebrew, das Binary auf macOS und Linux, der Wortlaut von `gh release download`.

**Bindung aus `slice-v1-abschluss-binary-namen`** (`AGENTS.md` §3.13): Der Slice setzt einen veröffentlichten stabilen Tag voraus, dessen Binaries `pgwire-recorder-<os>-<arch>` heißen. Einen solchen Tag bereitete `slice-v1-abschluss-release-v0-1-1` vor (Welle `welle-v1-nachtrag`, Reihenfolge 4: `VERSION`, Eintrag der Freigabe, Bezug im Handbuch); gesetzt hat ihn die Koordination (Claude) am 2026-10-21 nach dem Ja des Nutzers vom selben Tag („bitte auch den tag fahren“), nicht der Nutzer und nicht ein Implementer-Lauf (`SPEC-067` Punkt 18). **Bindung aus `slice-v1-abschluss-release-v0-1-1`** (`AGENTS.md` §3.13): Der Eintrag der Freigabe nennt das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot); Wiederanlauf des Jobs, Nachtrag von Hand und alle Belege des Taps gehören diesem Slice. Die Abschnitte „Binary vom GitHub-Release“ und „Container-Image aus der Registry“ in §2 *Installation* des Handbuchs und die Zeile zur Version in `README.md` schreibt `slice-v1-abschluss-release-v0-1-1`; dieser Slice ergänzt nur den Abschnitt zu Homebrew und die Plattform-Zeile (dieselbe Datei, andere Abschnitte, der Nachweis folgt dem Release-Slice).

**Entschieden vom Architect am 2026-10-21 (Randformen vor dem Code, `AGENTS.md` §3.12):** Der Nachweis unter Linux läuft als Werkzeug `make homebrew-nachweis` in einem Container des per Digest gepinnten Image `homebrew/brew`; die Schritte (`brew tap`, `brew install`, `pgwire-recorder version`, `brew test`, Start eines Replay-Laufs) liegen in `tools/harness/homebrew-nachweis.sh` und sind auf beiden Plattformen dieselben. Der Nachweis unter macOS ist die Handlung des Nutzers an seinem Rechner mit demselben Skript (`--schritte`). Was ein Beleg ist, was eine Angabe ohne Ausgabe wert ist, was „ein Replay-Lauf startet“ verlangt und welche Kombination aus Betriebssystem und Architektur zugesagt oder Grenze ist, steht als Vertragstext in `SPEC-067` Punkt 20 (`harness/vertraege/`) und in `LH-FA-19.a`; beides liegt vor dem ersten Code-Commit dieses Slice. **Nachgeführt am 2026-10-21 (§3.9):** Drei Teile von Punkt 20 hat der Architect erst **nach** dem ersten Lauf und nach dem Review entschieden — das Maß des Wartens vor dem `SIGTERM` samt Exit-Code 7 und die zulässigen Zeilen der Fehlerausgabe (Punkt 20 (3)), die Zusage des Kopfes über den Wächter (Punkt 20 (5)) und Lauf-Fassung und Kürzung des Belegs (Punkt 20 (6)); der Code dieses Slice folgt dieser Fassung, und §6 nennt die drei Randformen (8, 10, 11) in ihrer nachentschiedenen Form.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Erzeugen der Formel — `slice-v1-abschluss-homebrew`; hier wird das fertige Verfahren mit dem echten Release angewandt.
- Ein Gate oder eine Gegenprobe für `make homebrew-nachweis` — Schicht-Abgrenzung zur Gate-Kette: Das Ziel ruft Netz und den echten Tap, es urteilt nicht, und eine Gegenprobe mit Attrappen beliefe genau die Stelle, die der Nachweis belegen soll (`SPEC-067` Punkt 20 (4) und (7), Punkt 15).
- Eine Zeile in der Freigabe-Checkliste (`SPEC-067` Punkt 17) — es wäre ein anderer Vorgang; die Belege des Taps gehören dem Nachweis, nicht der Zeile des Jobs `homebrew` (Punkt 17 (8)).
- `brew audit` — Bestand bleibt bewusst stehen: Die Zeile `version` der Formel ist dort ein bekannter Befund (`SPEC-067` Punkt 19 (4)) und sperrt die Installation nicht.
- Die Kombinationen aus Betriebssystem und Architektur, für die kein Rechner bereitsteht — sie stehen als Grenze in `LH-FA-19.a`; ein emulierter Lauf trüge die Zusage nicht, und die Grenze ist entschieden, nicht offen (`AGENTS.md` §3.12, Plattform).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Der Tap trägt die Formel des ersten stabilen Tags mit dem Schema ohne Version (`v0.1.1`, `slice-v1-abschluss-release-v0-1-1`; Commit im Tap als Beleg; erzeugt vom Job `homebrew` des Tag-Laufs, bei Fehler per Wiederanlauf oder Nachtrag von Hand nach `SPEC-067` Punkt 19 (7)), und Abnahmeszenario 11 ist je Kombination aus Betriebssystem und Architektur belegt oder begrenzt (Form, Ablage und Grenze im Punkt darunter): belegt heißt Installation mit den dokumentierten Befehlen, `pgwire-recorder version` meldet die installierte Version, ein Replay-Lauf startet. Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung). Aus `slice-doku-ist-stand` (dort §6, *Installation*): Der Slice liefert die Befehle des Homebrew-Abschnitts in §2 *Installation* des Handbuchs, belegt mit der Installation aus dem echten Tap je Kombination, die er nennt, und die Plattform-Zeile in §1 *Voraussetzungen* (macOS).
- [ ] `SPEC-067` Punkt 20: `make homebrew-nachweis` liegt (Fragment unter `harness/mk/`, Skript unter `tools/harness/`, Zeile in der Werkzeug-Tabelle von `harness/README.md`), der Pin ist ein gegen die Registry aufgelöster Digest mit dem Tag als Kommentar, und ein Lauf gegen den echten Tap steht mit seiner Ausgabe in §7. Je Kombination aus Betriebssystem und Architektur liegt entweder der Beleg nach `SPEC-067` Punkt 20 (6) in §7 oder die Grenze in `LH-FA-19.a` und im Handbuch (`AGENTS.md` §3.11, §3.12).
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
| Tap `pt9912/homebrew-pgwire-recorder` | neu (extern) | Formel des ersten stabilen Tags mit dem neuen Schema, durch den Job `homebrew` (Wiederanlauf oder Nachtrag von Hand bei Fehler) |
| `docs/user/benutzerhandbuch.md` | update | Befehle gegen die echte Installation prüfen |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |
| `harness/vertraege/SPEC-067-release.md` | update (liegt) | Punkt 20: Werkzeug, Netz, Pin, Aufräumen, kein Gate, Beleg-Form, Replay-Lauf (Architect, vor dem Code) |
| `spec/spezifikation.md` (`LH-FA-19.a`) | update (liegt) | Zusage je Kombination, Grenze für die übrigen, was ein Nachweis ist (Architect, vor dem Code) |
| `tools/harness/homebrew-nachweis.sh` | neu | die Schritte des Nachweises, auf beiden Plattformen dieselben (`SPEC-067` Punkt 20 (2), (3)) |
| `harness/mk/homebrew-nachweis.mk` | neu | Ziel `homebrew-nachweis` mit `HOMEBREW_IMAGE ?=` (Digest-Pin), weder in `GATE_CHECKS` noch in einer Gegenproben-Liste |
| `harness/README.md` | update (liegt) | Zeile in der Tabelle der Werkzeuge (kein Gate) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-homebrew` ist `done` (das Verfahren der Formel; Reihenfolge in §5 von [welle-v1-nachtrag](../welle-v1-nachtrag.md), Entscheidung des Nutzers vom 2026-10-11), `slice-v1-abschluss-binary-namen` ist `done`, und **ein stabiler Tag nach beiden ist veröffentlicht**: Seine Binaries heißen `pgwire-recorder-<os>-<arch>`, und der Lauf hatte den Job `homebrew` (gesetzt hat den Tag die Koordination (Claude) nach dem Ja des Nutzers vom 2026-10-21, `SPEC-067` Punkt 18; `v0.1.0` zählt nicht, Entscheidung des Nutzers vom 2026-10-20). Der Slice, der diesen Tag vorbereitet und belegt, ist `slice-v1-abschluss-release-v0-1-1`; er liegt in `done/` (Start-Trigger, seit seiner Closure am 2026-10-21), und der Tag `v0.1.1` ist veröffentlicht (Lauf `39551891874`, Eintrag `docs/maintainer/freigaben/0.1.1.md`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Nachweis braucht eine zweite Plattform-Umgebung, die nicht bereitsteht — zurück zur Zerlegung.
- `in-progress` → `open`: Der Tap lässt sich nicht beschreiben (Zugang) — extern klären.
- `in-progress` → `next`: **zu groß, über 900** in der Metrik unten — zurück zur Zerlegung.

**Metrik der Marke** (`AGENTS.md` §3.14, Architect 2026-10-21): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness spec docs README.md`, Tests eingeschlossen. **Schätzung des Architect in dieser Metrik:** Skript `homebrew-nachweis.sh` 150 bis 260, Fragment 10 bis 20, Vertragstext in `SPEC-067` Punkt 20 und `LH-FA-19.a` (mit diesem Commit geliefert) rund 35, Zeile in `harness/README.md` 1, Handbuch-Abschnitt und `README.md` 40 bis 70, Plan 60 bis 110 — zusammen **rund 300 bis 500**. Der Faktor für Test- und Gegenprobe-Code entfällt: Das Ziel hat keine Gegenprobe und keinen Fall (`SPEC-067` Punkt 20 (7)). **Meldeschwelle 700, Marke 900** (über 900 greift die Rückführung oben); der Bericht des Implementers nennt die Messung in dieser Metrik.

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

**Randformen des Vertrags, entschieden vor dem Code** (`AGENTS.md` §3.12; Architect, 2026-10-21; der Vertragstext liegt in `SPEC-067` Punkt 20 (`harness/vertraege/`) und in `LH-FA-19.a`, beides vor dem ersten Code-Commit dieses Slice; die Randformen 8, 10 und 11 hat der Architect am 2026-10-21 nach dem Review im Vertrag nachentschieden — nach dem Lauf des Implementers, vor der Closure). Je Randform die Stelle, die sie entscheidet:

1. **Netz.** Der Lauf braucht Netz (Tap, Anhänge des GitHub-Release, Homebrew im Image), anders als jedes Gate dieser Kette; genau die Erreichbarkeit ist der Gegenstand des Nachweises — Punkt 20 (4).
2. **Kein Gate.** Das Ziel steht nicht in `GATE_CHECKS` und in keiner Gegenproben-Liste, kein Hook und kein Stand liest sein Ergebnis; es steht in der Tabelle der Werkzeuge in `harness/README.md` — Punkt 20 (4), (7).
3. **Pin.** Image `homebrew/brew:<Tag>@sha256:<64 Hexziffern>`, `HOMEBREW_IMAGE ?=` im Fragment, Tag als Kommentar (eine Version wie `4.6.20`, nicht `latest`, `main` oder `master`; Digest des Index, nicht eines Plattform-Manifests); der Implementer löst den Digest vor dem Einsetzen gegen die Registry auf und nennt Tag, Digest und Datum im Beleg; ein Platzhalter aus Nullen ist ungültig — Punkt 20 (1).
4. **Fehlendes Netz, nicht erreichbarer Tap.** Exit-Code des gescheiterten Schrittes, kein Überspringen, kein Grün ohne Lauf, kein Rückfall auf eine lokale Formel; der Nachweis gilt dann als nicht erbracht — Punkt 20 (4).
5. **Aufräumen.** `docker run --rm` (Container samt anonymen Volumes), Temp-Verzeichnis der Schritte entfällt bei Ende, `SIGINT` und `SIGTERM`, der Replay-Lauf wird dabei beendet; nichts im Arbeitsbaum — Punkt 20 (1), (5).
6. **Aufrufarten des Treibers.** `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` haben keinen Fall, weil kein Hook und kein Stand das Ergebnis liest (Grund wie `SPEC-067` Punkt 13) — Punkt 20 (5).
7. **Keine Interaktion.** `docker run` ohne `-t`, Standardeingang `/dev/null`, `HOMEBREW_NO_ANALYTICS`, `HOMEBREW_NO_AUTO_UPDATE`, `HOMEBREW_NO_ENV_HINTS` gesetzt; weitere `HOMEBREW_*`-Variablen nur, wenn keine eine Summe, eine Adresse oder den Tap ersetzt, und im Beleg genannt — Punkt 20 (1).
8. **Ausführungsort des Skripts.** `make homebrew-nachweis` unter Linux mit Docker; `--schritte` am Host ist allein die Handlung des Betreibers (macOS: bash 3.2, deshalb ohne `mapfile`, ohne assoziative Arrays, ohne `${v^^}`), und ein Agentenlauf fährt sie nicht (`AGENTS.md` §3.1; sein Weg ist `make homebrew-nachweis` im Container). Das ist eine **Regel an den Leser, kein Mechanismus**: Der Wächter `pretooluse-command-guard.sh` blockt `brew install …` am Kopf eines Segments, liest aber den Inhalt einer Datei nicht und lässt `bash tools/harness/homebrew-nachweis.sh --schritte` durch (am 2026-10-21 an beiden Eingaben gemessen); sein eigener Kopf nennt ihn einen Stolperdraht und keine Sandbox. Der Kopf des Skripts sagt deshalb, wem der Aufruf am Host gehört und dass ein Agentenlauf ihn nicht fährt, und sagt dem Wächter keinen Schutz zu, den er nicht leistet. **Akzeptiertes Negativ, mit Grund im Vertrag:** Der Aufruf am Host wird nicht gesperrt — unter macOS ist er die vorgesehene Form des Nachweises, und ein Muster im Wächter auf diesen Dateinamen deckte weder eine Kopie noch `sh <Datei>` und erzeugte den Anschein einer Deckung; der Wächter wird dafür nicht geändert — Punkt 20 (5).
9. **Plattform der Zusage** (`AGENTS.md` §3.12, Plattform; `LH-QA-03`). Die Zusage gilt je Kombination aus Betriebssystem und Architektur einzeln: zugesagt ist, was mit seinen Ausgaben belegt ist; jede übrige der vier Kombinationen (macOS und Linux, je `amd64` und `arm64`) ist eine **Grenze** in `LH-FA-19.a`, und das Handbuch nennt nur die geprüften (`AGENTS.md` §3.11). Windows bleibt außerhalb (`SPEC-042`). Die Grenze ist damit entschieden, nicht offen — sie wird kein Risiko mit Ausgang *weiter offen* (`BEO-REPO/verhalten-nur-unter-linux-geprueft`).
10. **Beleg-Form und Ablage.** Je Kombination: Datum, `uname -s`, `uname -m`, unter Linux Tag und Digest des Image, die Befehle (a) bis (e) mit ihren Ausgaben, die Version und der Exit-Code des Replay-Laufs; Ablage in §7 dieses Plans (Form wie `SPEC-067` Punkt 18). Eine Angabe ohne diese Ausgaben — auch die des Nutzers — ist eine Angabe der Koordination, kein Beleg: zitiert mit Datum und Wortlaut, die Kombination gilt als nicht geprüft. Der Beleg stammt aus einem Lauf des Skripts in der Fassung, die der Slice liefert; ändert sich das Skript danach an einer Stelle, die seine Ausgabe berührt, zählt der Beleg des neuen Laufs. **Gekürzt** wird allein der Teil einer Ausgabe, der ein Geheimnis oder einen kurzlebigen Token trägt (Signatur, JWT und Ablaufzeit einer signierten Adresse, ein Token, ein Schlüssel, ein Passwort), genau dieser Teil, mit der Marke `[gekürzt: <Grund>]` im Text und einem Satz im Beleg, welche Zeilen gekürzt sind und warum, und beim Einfügen vor dem Commit; eine Kürzung aus einem anderen Grund oder ohne Marke ist ein Beleg, der die Ausgabe nicht vollständig zeigt — Punkt 20 (6).
11. **Was „ein Replay-Lauf startet“ verlangt.** `pgwire-recorder replay --listen 127.0.0.1:55432 --input <Datei>` mit der Aufzeichnung aus Punkt 20 (3): **verlangt** sind mindestens eine Zeile `level=INFO` und Exit-Code 0 nach `SIGTERM` bzw. `Strg+C`; **zulässig** sind daneben allein weitere Zeilen `level=INFO` und genau eine Zeile `level=WARN` mit `code=PGR-W2001` (nie zugeordnete Session — erwartet, weil keine Anwendung verbindet); **jede andere Zeile** (`level=ERROR`, eine Warnung mit einem anderen Code, eine zweite Zeile mit `PGR-W2001`, eine Zeile ohne Feld `level=`) lässt den Nachweis nicht durchgehen: Der Beleg nennt sie, der Nachweis gilt als nicht erbracht, und der Lauf geht als Randform-Rückgabe an den Architect. Keine Verbindung einer Anwendung, kein Server. **Maß des Wartens vor dem `SIGTERM`:** ein Ereignis mit Grenze statt eines festen Maßes — das Skript wartet, bis die Datei der Fehlerausgabe nicht mehr leer ist (allein ihre Größe, nicht ihr Inhalt; die eine Ausnahme von Punkt 20 (2)), im Abstand von höchstens einer Sekunde, Grenze 60 Sekunden; läuft sie ab, ist das Warten der gescheiterte Schritt, das Skript endet mit Exit-Code 7, und der Nachweis gilt als nicht erbracht (`BEO-REPO/test-faellt-unter-cpu-last-aus`). Das Warten auf das Ende nach dem `SIGTERM` bleibt dagegen ohne Grenze — entschieden, nicht offen: Dass der Lauf nach dem Signal endet, ist Gegenstand des Nachweises, und ein Hängen ist selbst der Befund. Auf beiden Plattformen dieselbe Form — `LH-FA-19.a` und Punkt 20 (3).
12. **Version des Nachweises.** Der Tap trägt die höchste stabile Version; der Beleg nennt die Version, die er installiert hat, und gilt für sie — Punkt 20 (6), Punkt 19 (5).
13. **Kein Fall, keine Gegenprobe, keine Abdeckungs-Deklaration.** Das Ziel urteilt nicht und sagt nichts zu, was ein Gate oder Hook glaubt; `AGENTS.md` §3.10 greift deshalb nicht, und eine Gegenprobe mit Attrappen beliefe genau die Stelle, die der Nachweis belegen soll — Punkt 20 (7). **Akzeptiertes Negativ, mit Grund:** Ein Fehler im Skript selbst (ein Schritt, der versehentlich nicht läuft) fiele nur dem Leser der Ausgabe auf; getragen wird das vom Review, der jede Zeile des Handbuchs gegen den Beleg in §7 prüft.

- Das Verhalten von Homebrew gegenüber Drittanbieter-Taps (`brew trust`) ist erst mit dem echten Tap belegbar — **Ausgang:** offen bis Closure.
- Ob `brew install` einer Formel ohne Archiv im Image `homebrew/brew` als dessen Benutzer gelingt (Rechte, vorhandenes `git`, Verhalten ohne TTY), zeigt erst der Lauf. Scheitert ein Schritt an der Umgebung des Image und nicht an der Formel, hält der Implementer an und gibt Meldung und Schritt dem Architect zurück (Vertragstext, kein Code: `SPEC-067` Punkt 20 (1) oder (2)) — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-release-v0-1-1` (Bindung, §1): Ergebnis des Jobs `homebrew` im Lauf des Tags — grün: der Commit im Tap ist der Beleg des ersten Liefer-Punkts; rot: Wiederanlauf nur dieses Jobs oder Nachtrag von Hand, kein neuer Tag (`SPEC-067` Punkt 19 (7)); die Handbuch-Abschnitte zu Binary und Image schreibt der Release-Slice — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-release-v0-1-1` (§6 Risiko (a), Ausgang *eingetreten*, bei dessen Closure am 2026-10-21 eingetragen, `AGENTS.md` §3.13): Die erste Außenwirkung des Jobs `homebrew` ist gelaufen und hat einen Teil der Annahmen widerlegt. **Gemessener Stand**, den dieser Slice nicht erneut annimmt, sondern gegen den Stand seines eigenen Laufs misst: Der Job lief im Lauf `39551891874` viermal — Versuche 1 bis 3 rot mit `release: homebrew: tap-push: der Tap weist den Push ab (Exit 128)` (`Permission … denied`, `error: 403`), Versuch 4 grün mit `release: homebrew: Formel 0.1.1 im Tap abgelegt`; die Spitze von `main` im Tap ist seither `6c4ac72843996939b4295b01aeb7cbe9e759a7ac` (2026-10-21 06:14:17Z). Das Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` wurde dabei nicht neu gesetzt (`gh secret list`, `updated_at` 2026-10-20T07:41:39Z); *welches* Recht der Betreiber zwischen den Versuchen änderte, zeigen weder Repo noch Lauf (Angabe der Koordination, `SPEC-067` Punkt 15). **Weiter unbelegt und damit Gegenstand dieses Slice:** die Annahme der Formel durch Homebrew, `brew tap` und `brew install` auf macOS und Linux, der Wortlaut von `gh release download` mit dem echten `gh` und das Schreibrecht des Geheimnisses ohne Eingriff des Betreibers — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-homebrew` (§6 dort), je mit Ausgang bei Closure dieses Slice: (a) das Schreibrecht des Tap-Tokens — das Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` (Fine-grained PAT nur für den Tap, Contents: Read and write, Entscheidung des Nutzers vom 2026-10-20) zeigt der Job `homebrew` des Laufs, der die erste Formel erzeugt; läuft die Formel stattdessen per Nachtrag mit dem eigenen Token des Nutzers, ist das Recht des Geheimnisses nicht belegt, und der Ausgang ist *weiter offen* (Register); (b) die Annahme der Formel durch Homebrew (Syntax, die Zeile `version`, Installation einer Datei ohne Archiv, `test do`) — eine Korrektur geht über die Vorlage und einen erneuten Nachtrag, ohne neuen Tag; (c) der Wortlaut von `gh release download` mit dem echten `gh` — **Ausgang:** offen bis Closure.
- Aus `slice-v1-abschluss-homebrew` (Verifikation V-1, Closure dort): Ob „Re-run failed jobs“ bei GitHub nur den Job `homebrew` wiederholt, ohne das Image neu zu bauen, zeigt kein Test dieses Repos; `docs/maintainer/releasing.md` sagt es nicht mehr zu. Muss dieser Slice den Job wiederholen, nennt sein Beleg, was GitHub tat, und `releasing.md` §6 und §11 bekommen den Satz erst dann — **Ausgang:** offen bis Closure.

**Stand der Randformen und der Risiken nach dem Lauf** (Implementer, 2026-10-21, §3.9; der Text der Randformen oben bleibt unverändert, hier steht allein, was der Lauf zeigte — die Ausgänge der Risiken setzt die Closure):

- **Randform 3 (Pin):** umgesetzt als `homebrew/brew:4.6.20@sha256:b0072bfd…fdedfd`; der Digest ist am 2026-10-21 gegen die Registry aufgelöst und ist der des Index (`MediaType: application/vnd.oci.image.index.v1+json`, Manifeste `linux/amd64` und `linux/arm64`), Nachweis in §7. Kein Platzhalter.
- **Randform 5 (Aufräumen):** nach dem Lauf kein Container (`docker ps -a` leer), kein Volume (`docker volume ls -q | wc -l` gibt `0`), nichts Neues im Arbeitsbaum (`git status --short` leer).
- **Randform 10 (Beleg-Form):** Das Skript gibt `date -u`, `uname -s` und `uname -m` als erste Befehle aus, damit die Angaben des Belegs aus dem Lauf selbst stammen und auf beiden Plattformen dieselbe Form haben; Tag und Digest nennt die erste Zeile des Rezepts. Das ist die Umsetzung von Punkt 20 (6), keine Entscheidung einer Randform.
- **Randform 11 (Replay-Lauf):** erfüllt — drei Zeilen `level=INFO`, keine `level=ERROR`, Exit-Code `0` nach `SIGTERM`, und genau eine Zeile `level=WARN` mit `code=PGR-W2001`: die eine zulässige und in diesem Szenario erwartete Warnung nach Punkt 20 (3), weil keine Anwendung verbindet. Keine weitere Zeile. Das Warten vor dem `SIGTERM` ist seit der Nachentscheidung des Architect ein Ereignis mit Grenze (Größe der Fehlerausgabe, Abstand 1 s, Grenze 60 s, sonst Exit-Code 7); beide Wege sind gefahren (§7, Nacharbeit zu F-6).
- **Risiko `brew trust`:** im Lauf nicht eingetreten — Homebrew hat für den Drittanbieter-Tap keine Bestätigung verlangt (`brew tap` meldet `Tapped 1 formula`, `brew install` läuft ohne Rückfrage).
- **Risiko Umgebung des Image:** im Lauf nicht eingetreten — `brew tap` klont mit dem im Image vorhandenen `git` als Benutzer `linuxbrew`, `brew install` legt die Datei ohne Archiv an (`4 files, 11.4MB`), `brew test` läuft ohne TTY. Kein Halt, keine Rückgabe an den Architect.
- **Was dieser Lauf nicht berührt:** der Wortlaut von `gh release download` mit dem echten `gh` (das ist der Job `homebrew` auf dem Runner, den dieser Slice nicht erneut gefahren hat), das Schreibrecht des Geheimnisses ohne Eingriff des Betreibers und der Wiederanlauf eines Jobs (der Job war im vierten Versuch grün, ein Wiederanlauf war nicht nötig).

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

### Angaben der Koordination — kein Beleg (vor dem Start eingetragen, Planner, 2026-10-21)

Eingaben für den Implementer, keine Belege: Keine der beiden Angaben trägt die
Ausgaben, die `SPEC-067` Punkt 20 (6) für einen Nachweis verlangt (Datum,
`uname -s`, `uname -m`, die Befehle (a) bis (e) mit ihren Ausgaben, Version und
Exit-Code). Beide stehen hier mit Datum und Wortlaut, wie Randform 10 in §6 es
vorschreibt; die betroffene Kombination gilt als **nicht geprüft**, bis ein
Nachweis in dieser Sektion liegt.

1. **macOS, Architektur unbekannt (Nutzer, 2026-10-21).** Der Nutzer sagte der
   Koordination: „ich habe auf macos eben den test erfolgreich abgeschlossen
   (brew install und running)“. Mehr ist nicht erfragt und nicht bekannt: weder
   die Architektur (`uname -m`) noch die Ausgabe von `pgwire-recorder version`
   noch ein Replay-Start nach `SPEC-067` Punkt 20 (3). Die Angabe sagt damit,
   dass ein Nachweis auf macOS aussichtsreich ist — sie ersetzt ihn nicht. Beide
   macOS-Kombinationen stehen bis zu einem Nachweis mit Ausgaben als Grenze in
   `LH-FA-19.a` und im Handbuch (`AGENTS.md` §3.11, §3.12).
2. **Keine Option `--version` (Nutzer, 2026-10-21).** Der Nutzer beobachtete,
   dass das Binary keine Option `--version` kennt. Das ist der Ist-Zustand und
   kein Befund gegen diesen Slice: Der Nachweis ruft `pgwire-recorder version`
   als Kommando (`SPEC-067` Punkt 20 (2) (c)), nicht als Option, und Handbuch
   wie Beleg nennen das Kommando. Den Gegenstand der Beobachtung führt
   `slice-version-option` in seinem §1; dieser Slice weist ihm nichts zu, hängt
   nicht von ihm ab und ändert die Stellen nicht, die er dort nennt. Liegt
   `slice-version-option` beim Nachweis-Lauf schon in `done/`, bleibt das
   Kommando `version` dennoch der Weg dieses Slice (dort §1, Klasse 2: das
   Kommando bleibt unverändert).

### Belege des Implementers (2026-10-21, Zweig `impl-nachweis`, Worktree `wt-nachweis`, Dev-Container `pgr-dev2`)

#### Nachweis Linux · x86_64 (Abnahmeszenario 11, Form nach `SPEC-067` Punkt 20 (6))

> **Dieser Beleg stammt aus einem Lauf der gelieferten Fassung** (Punkt 20 (6)): Stand `74e19478`, also nach der Nacharbeit zu F-1, F-2 und F-6. Zwei frühere Läufe sind damit ersetzt und zählen nicht als Beleg — der erste (08:38:05Z) stammt aus der Fassung ohne die beiden Exit-Code-Zeilen, der zweite (09:17:20Z) aus der Fassung mit `sleep 5`; beide endeten mit Exit 0 und liegen als `/tmp/nachweis-lauf1.log` und `/tmp/nachweis-lauf2.log` im Dev-Container.
>
> **Gekürzt** ist allein die Zeile `==> Downloading from https://release-assets.githubusercontent.com/…`, und zwar genau der Teil, der Signatur, JWT und Ablaufzeit der signierten Adresse trägt (kurzlebiger Token, Ablauf `2026-10-21T10:04:42Z`); Dienst, Pfad und Dateiname bleiben lesbar, die Stelle trägt die Marke `[gekürzt: …]`. Keine andere Zeile ist gekürzt (Punkt 20 (6)).

- **Datum:** `2026-10-21T09:29:42Z` (Ausgabe von `date -u` im Container des Laufs).
- **`uname -s`:** `Linux` · **`uname -m`:** `x86_64`. Dazu, als Beleg der Plattform des Docker-Servers (F-5): `docker version --format '{{.Server.Os}}/{{.Server.Arch}}'` gibt `linux/amd64` (Client ebenso), `docker info --format '{{.OSType}}/{{.Architecture}}'` gibt `linux/x86_64`. Die Kombination des Belegs ist damit **Linux/amd64**.
- **Image:** Tag `4.6.20`, Digest `sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd`. Der Digest ist am 2026-10-21 gegen die Registry aufgelöst worden (`docker buildx imagetools inspect homebrew/brew:4.6.20`: `MediaType: application/vnd.oci.image.index.v1+json`, Manifeste für `linux/amd64` und `linux/arm64`; `--raw | sha256sum` gibt denselben Wert) — also der Digest des **Index**, nicht eines Plattform-Manifests.
- **Umgebung:** `HOMEBREW_NO_ANALYTICS=1`, `HOMEBREW_NO_AUTO_UPDATE=1`, `HOMEBREW_NO_ENV_HINTS=1`; keine weitere `HOMEBREW_*`-Variable. `docker run --rm` ohne `-t`, Standardeingang `/dev/null`, allein `tools/harness/homebrew-nachweis.sh` schreibgeschützt eingehängt.
- **Befehl:** `make homebrew-nachweis` · **Dauer:** 14,4 s · **Exit-Code des Skripts:** `0` (die letzte Zeile des Skripts nennt ihn selbst).
- **Version aus (c):** `pgwire-recorder 0.1.1` · **Exit-Code des Replay-Laufs aus (3):** `0`.

Ausgabe des Laufs, vollständig bis auf die eine markierte Kürzung oben:

```text
homebrew-nachweis: Image homebrew/brew:4.6.20@sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd
+ date -u +%Y-%m-%dT%H:%M:%SZ
2026-10-21T09:29:42Z
+ uname -s
Linux
+ uname -m
x86_64
+ brew tap pt9912/pgwire-recorder
==> Tapping pt9912/pgwire-recorder
Cloning into '/home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/pt9912/homebrew-pgwire-recorder'...
Tapped 1 formula (13 files, 68.6KB).
+ brew install pt9912/pgwire-recorder/pgwire-recorder
==> Fetching downloads for: pgwire-recorder
==> Fetching pt9912/pgwire-recorder/pgwire-recorder
==> Downloading https://github.com/pt9912/pgwire-recorder/releases/download/v0.1.1/pgwire-recorder-linux-amd64
==> Downloading from https://release-assets.githubusercontent.com/github-production-release-asset/1403356322/e5be8986-543a-4a75-8c54-5946309ffb16?rscd=attachment%3B+filename%3Dpgwire-recorder-linux-amd64&rsct=application%2Foctet-stream&… [gekürzt: Signatur, JWT und Ablaufzeit der signierten Adresse]
==> Installing pgwire-recorder from pt9912/pgwire-recorder
🍺  /home/linuxbrew/.linuxbrew/Cellar/pgwire-recorder/0.1.1: 4 files, 11.4MB, built in 3 seconds
==> Running `brew cleanup pgwire-recorder`...
+ pgwire-recorder version
pgwire-recorder 0.1.1
+ brew test pgwire-recorder
Warning: test is a developer command, so Homebrew's
developer mode has been automatically turned on.
To turn developer mode off, run:
  brew developer off

✔︎ JSON API cask_tap_migrations.jws.json
==> Testing pt9912/pgwire-recorder/pgwire-recorder
==> /home/linuxbrew/.linuxbrew/Cellar/pgwire-recorder/0.1.1/bin/pgwire-recorder version
+ pgwire-recorder replay --listen 127.0.0.1:55432 --input /tmp/tmp.4bIPcMQH9K/aufzeichnung.yaml (Fehlerausgabe in eine Datei)
+ warten, bis der Replay-Lauf seine erste Zeile geschrieben hat (Groesse der Fehlerausgabe, Abstand 1 s, Grenze 60 s)
+ kill -TERM <Replay-Lauf>
--- Fehlerausgabe des Replay-Laufs ---
time=2026-10-21T09:29:55.513Z level=INFO msg="replay gestartet" listen=127.0.0.1:55432 input=/tmp/tmp.4bIPcMQH9K/aufzeichnung.yaml
time=2026-10-21T09:29:56.511Z level=INFO msg="Herunterfahren begonnen" sessions=0
time=2026-10-21T09:29:56.511Z level=WARN msg="1 aufgezeichnete Session(s) nie zugeordnet, die erste mit Kennung 1" code=PGR-W2001
time=2026-10-21T09:29:56.511Z level=INFO msg="replay beendet"
--- Exit-Code des Replay-Laufs: 0
homebrew-nachweis: Skript endet mit Exit-Code 0
homebrew-nachweis: der Lauf endet mit Exit-Code 0
```

**Was der Lauf zeigt, Schritt für Schritt.** (a) `brew tap pt9912/pgwire-recorder` ohne zweites Argument leitet die Adresse auf `pt9912/homebrew-pgwire-recorder` ab (Zeile `Cloning into '…/Taps/pt9912/homebrew-pgwire-recorder'`) und nimmt die eine Formel an (`Tapped 1 formula`) — eine Bestätigung wie `brew trust` hat Homebrew nicht verlangt. (b) `brew install` lädt das Binary des Tags `v0.1.1` vom GitHub-Release unter dem Namen ohne Version (`pgwire-recorder-linux-amd64`) und installiert es ohne Abbruch; Homebrew prüft die Summe der Formel beim Laden, eine Abweichung bräche den Schritt mit `SHA256 mismatch`. (c) `pgwire-recorder version` gibt `pgwire-recorder 0.1.1` — die Version der Formel im Tap. (d) `brew test pgwire-recorder` fährt den Block `test do` der Formel (`==> /home/linuxbrew/.linuxbrew/Cellar/pgwire-recorder/0.1.1/bin/pgwire-recorder version`) ohne Befund. (e) Der Replay-Lauf gegen die im Skript erzeugte Aufzeichnung schreibt drei Zeilen `level=INFO` und endet nach `SIGTERM` mit Exit-Code `0`; dazu steht genau eine Zeile `level=WARN` mit `code=PGR-W2001` (nie zugeordnete Session) — die eine zulässige und in diesem Szenario erwartete Warnung nach `SPEC-067` Punkt 20 (3), weil keine Anwendung verbindet. Keine Zeile `level=ERROR`, keine Warnung mit einem anderen Code, keine zweite mit `PGR-W2001`, keine Zeile ohne Feld `level=`; damit trägt die Fehlerausgabe genau die Zeilen, die Punkt 20 (3) verlangt und zulässt. Vor dem `SIGTERM` steht die Zeile des Wartens: Das Skript hat gewartet, bis die Fehlerausgabe nicht mehr leer war (eine Sekunde, weit unter der Grenze von 60 s).

**Was dieser Beleg nicht trägt.** Er gilt für **Linux/amd64** und für die Version `0.1.1`. `linux/arm64` ist nicht gefahren (kein Rechner, kein emulierter Lauf, §1 Abgrenzung); macOS ist nicht gefahren (die Angabe des Nutzers oben trägt keine Ausgaben und bleibt eine Angabe der Koordination). Diese drei Kombinationen stehen als **Grenze** in `LH-FA-19.a` und im Handbuch. `brew audit` gehört nicht zum Nachweis (§1 Abgrenzung). Ein Fehler im Skript selbst fiele nur dem Leser der Ausgabe auf (§6 Randform 13, akzeptiertes Negativ).

#### Zusagen des Werkzeugs und ihre Gegenprobe

`SPEC-067` Punkt 20 (7) schließt Fälle, Gegenprobe und Abdeckungs-Deklaration für dieses Ziel aus; `AGENTS.md` §3.10 greift damit nicht (§6 Randform 13). Was statt einer Gegenprobe belegt ist, steht hier:

| Zusage (Punkt 20) | Beleg aus dem Lauf | Gegenbeispiel |
|---|---|---|
| (1) Pin ist Tag + Digest des Index, kein Platzhalter | Zeile 1 der Ausgabe nennt beide; `docker buildx imagetools inspect` nennt `MediaType: …image.index.v1+json` und beide Plattform-Manifeste | nicht gefahren — ausgeschlossen (Punkt 20 (7)); ein Platzhalter aus Nullen wäre an der Zeile erkennbar, ein Plattform-Digest am `MediaType` des `inspect` |
| (2) Reihenfolge (a) bis (e), jeder Befehl vor seiner Ausführung auf stdout, kein `brew audit` | die Zeilen `+ …` in genau dieser Reihenfolge, keine Zeile `+ brew audit` | nicht gefahren — ausgeschlossen (Punkt 20 (7)) |
| (2) Exit-Code: der des ersten scheiternden Schrittes, sonst 0 | die beiden letzten Zeilen (`Skript endet mit Exit-Code 0`, `der Lauf endet mit Exit-Code 0`); über `make` sähe der Aufrufer bei einem roten Lauf nur die `2` von make (`BEO-REPO/vertrag-nennt-exit-code-des-skripts-nicht-von-make`) | **gefahren, drei Fehlerpfade** — siehe die Nacharbeit zu F-1 unten: nicht ziehbares Image (Exit 125, keine Skript-Zeile), Lauf ohne Netz (Skript-Zeile mit Exit 1), grüner Lauf (beide Zeilen mit 0) |
| (3) Replay-Lauf startet: `level=INFO`, kein `level=ERROR`, Exit 0 nach `SIGTERM` | die vier Zeilen der Fehlerausgabe und `Exit-Code des Replay-Laufs: 0` | nicht gefahren — ausgeschlossen (Punkt 20 (7)); die Zusage ist die Ausgabe selbst, die ein Mensch liest |
| (4) kein Gate | `make -pn` nach dem Lauf: `GATE_CHECKS` führt das Ziel nicht (13 Einträge, `homebrew-nachweis` nicht darunter), `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST` ebenso wenig; `grep -rln homebrew-nachweis --exclude-dir=.git .` nennt unter `harness/mk/`, `tools/` und `.claude/` allein das eigene Fragment und das eigene Skript — kein Hook, kein anderes Fragment | nicht gefahren — ausgeschlossen (Punkt 20 (7)) |
| (5) Aufräumen: `--rm`, Temp-Verzeichnis entfällt, nichts im Arbeitsbaum | nach dem Lauf: `docker ps -a --filter ancestor=homebrew/brew` leer, `docker ps -a` ohne jeden Container, `docker volume ls -q \| wc -l` gibt `0`, `git status --short` leer (die beiden Dateien des Werkzeugs lagen da schon im Commit `a2fce55f`) | nicht gefahren — ausgeschlossen (Punkt 20 (7)) |

#### Läufe des Implementers (Zweig `impl-nachweis`, Dev-Container `pgr-dev2`, 2026-10-21)

| Lauf | Stand | Ergebnis | Dauer |
|---|---|---|---|
| `make homebrew-nachweis` (Beleg oben) | `74e19478` | grün, Exit-Code des Skripts `0` | 14,4 s |
| Sonde zur abgelaufenen Grenze (F-6) | `74e19478` | Exit `7` wie zugesagt, obwohl der Lauf `0` lieferte | 3 s |
| `make homebrew-nachweis HOMEBREW_IMAGE=<nicht ziehbar>` (F-1) | `8632a78f` | Exit `125`, keine Skript-Zeile | 2 s |
| `make homebrew-nachweis HOMEBREW_IMAGE="--network none <Pin>"` (F-1) | `8632a78f` | Skript-Zeile mit Exit `1` (`brew tap` ohne Netz) | 2 s |
| `make test-integration` allein | `b39beea7` | grün, Exit 0; `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` `PASS` in 19,85 s | 2 min 20 s |
| `make gates` (erster Lauf) | `e81a303c` | grün, Exit 0 (alle 13 Einträge von `GATE_CHECKS`), `Serverfehler_ohne_Vorbild` `PASS` in 15,28 s | 4 min 42 s |
| `make gates` (nach der ersten Nacharbeit) | `b39beea7` | **rot, Exit 2** — `test-integration`: `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` `FAIL` nach 30,02 s mit `DROP DATABASE IF EXISTS vgl_abw: timeout: context deadline exceeded`, Last 6,0 auf 20 Kernen; der Unterfall allein danach grün (Zeile darüber) — das ist `BEO-REPO/test-faellt-unter-cpu-last-aus`, **fünftes Auftreten**, kein Befund gegen diesen Slice | 4 min 2 s |
| `make gates`, `make docs-check` (Endstand) | siehe Zeile unten | — | — |
| `make docs-check` | `e81a303c`, `a9bd85f6` | grün, 1322 Dateien, 0 Befunde | je ~2 s |

Keine Gegenprobe gefahren: Dieses Ziel hat keine (Punkt 20 (7)), und für die unveränderten Gates sieht der Lauf keine Änderung (Entscheidung des Nutzers vom 2026-10-21).

**Größenmarke** (§4, Metrik `git diff --numstat f8bc8380..<letzter Commit> -- tools harness spec docs README.md`): am Stand `74e19478` Einfügungen 417, Löschungen 13, **Summe 430** (ohne die Zeilen dieses Nachtrags selbst — ein Plan kann die Messung seines eigenen Commits nicht enthalten; mit ihnen bleibt der Wert unter 500). Darin stecken 69 Zeilen Review-Report und 23 Zeilen Vertragstext des Architect (`SPEC-067`, `spec/spezifikation.md`); auf die Arbeit des Implementers entfallen 338 (Skript 129, Plan 139, Fragment 38, Handbuch 27, `README.md` 5). Das liegt unter der Meldeschwelle 700 und weit unter der Marke 900; die Schätzung des Architect (300 bis 500) trifft zu.

#### Was der Implementer beim Schreiben entschieden hat (kein Vertragstext)

Punkt 20 lässt diese vier offen; keine davon ist eine Randform im Sinne von `AGENTS.md` §3.12 (keine entscheidet über Eingaben, Formen oder Fehlerwege eines Vertrags), und keine ändert eine Zusage:

1. **Wartezeit vor dem `SIGTERM`:** fest 5 Sekunden. Punkt 20 (3) sagt „wartet“ ohne Maß; ein Warten auf eine Zeile der Ausgabe wäre eine Auswertung und ist durch Punkt 20 (2) ausgeschlossen.
2. **Ort des eingehängten Skripts im Container:** `/tmp/homebrew-nachweis.sh`, schreibgeschützt; das Temp-Verzeichnis der Schritte legt `mktemp -d` daneben.
3. **Kopf der Ausgabe:** `date -u`, `uname -s`, `uname -m` als erste Befehle (§6, Stand der Randformen).
4. **Kürzung im Beleg oben:** allein die zweite `Downloading`-Zeile, hinter dem Pfad; gekürzt sind Signatur und JWT der zeitlich begrenzten Adresse, die Kürzung steht im Text. Sonst ist die Ausgabe vollständig.
5. **`$(strip $(HOMEBREW_IMAGE))` an beiden Verwendungsstellen des Rezepts** (nachgetragen zu F-9): Punkt 20 (1) verlangt den Tag als Kommentar hinter der Pin-Zeile; `HOMEBREW_IMAGE ?= <Pin> # 4.6.20 …` lässt den Wert in `make` auf ein Leerzeichen enden, und die Docker-Referenz trüge es mit. `harness/mk/workflows.mk` löst dasselbe ohne Leerzeichen vor dem `#` und braucht kein `strip`; hier bleibt das Leerzeichen, weil der Pin mit Digest schon lang ist und `… # 4.6.20, Digest des Index, …` lesbarer bleibt. Beide Wege erfüllen Punkt 20 (1).

#### Nacharbeit zum Review vom 2026-10-21 (`docs/reviews/2026-10-21-review-slice-erster-release-homebrew-nachweis.md`, 0 HIGH, 3 MEDIUM, 3 LOW, 3 INFO)

Dem Implementer gehören F-1, F-4, F-5, F-7 und F-9; F-2, F-3, F-6 und F-8 sind Vertragstext und gehören dem Architect (`SPEC-067` Punkt 20 (3) und (5), `spec/spezifikation.md`, Plan §6) — dieser Nachtrag fasst §6 und den Vertragstext nicht an.

**F-1 (MEDIUM, behoben, Commit `8632a78f`).** Die alte Rezeptzeile schrieb jeden Code von `docker run` dem Skript zu. Jetzt gilt: Das **Skript** nennt seinen eigenen Code als letzte Zeile (`homebrew-nachweis: Skript endet mit Exit-Code <n>`, aus der `EXIT`-Falle, die den Code nicht verändert), und das **Rezept** nennt den Code des **Laufs** mit vier Fällen — `0`; `7`, den allein das Skript für die abgelaufene Grenze des Wartens wählt (seit F-6, außerhalb der Codes 0 bis 6 des Produkts und keiner der CLI-Codes); die CLI-Codes `125|126|127` mit dem Hinweis auf Punkt 20 (4); und sonst der Verweis auf die Skript-Zeile darüber. Kommentarblock des Fragments und Skript-Kopf sagen dasselbe. Drei Läufe, je mit der rot gesehenen Unterscheidung:

| Lauf | Ergebnis | Was er zeigt |
|---|---|---|
| `make homebrew-nachweis HOMEBREW_IMAGE=pgr-review-nicht-vorhanden:kein-tag` | `docker: Error response from daemon: pull access denied …`; `homebrew-nachweis: der Lauf endet mit Exit-Code 125; mit 125, 126 und 127 meldet die Docker-CLI …`; **keine** Zeile `Skript endet mit Exit-Code`; `make: *** … Error 125`, `make`-Exit `2` | der Fall des Reviewers: der Code stammt von der CLI, und die Ausgabe behauptet das Gegenteil nicht mehr |
| `make homebrew-nachweis HOMEBREW_IMAGE="--network none <Pin>"` | `brew tap` scheitert mit `Could not resolve host: github.com` (Punkt 20 (4), fehlendes Netz); `homebrew-nachweis: Skript endet mit Exit-Code 1`; `der Lauf endet mit Exit-Code 1; wessen Code das ist, sagt die Zeile … darüber`; `make`-Exit `2` | das Gegenstück: das Skript lief, der Code ist seiner, und die Zeile belegt es |
| `make homebrew-nachweis` | beide Zeilen mit `0`, 14,4 s (der Beleg oben) | der grüne Fall, Exit-Code 0 vom Skript |
| Sonde zur abgelaufenen Grenze (F-6, Zeile darunter) | `Skript endet mit Exit-Code 7`, `der Lauf endet mit Exit-Code 7; den waehlt das Skript …` | der Code `7` wird nicht als Code der Docker-CLI gelesen |

**F-6 (LOW, behoben, Commit `74e19478`; Vertragstext des Architect in Punkt 20 (2), (3)).** Das feste `sleep 5` ist weg. Das Skript nennt das Warten wie einen Schritt (`+ warten, bis der Replay-Lauf seine erste Zeile geschrieben hat (Groesse der Fehlerausgabe, Abstand 1 s, Grenze 60 s)`) und wartet, bis `[ -s <Fehlerausgabe> ]` gilt — allein die Größe, nie der Inhalt —, höchstens 60 Durchläufe im Abstand von einer Sekunde. Läuft die Grenze ab, nennt es sie, schickt wie sonst `SIGTERM`, gibt Fehlerausgabe und Exit-Code des Laufs aus und endet mit `7`. Das `wait` nach dem `SIGTERM` bleibt ohne Grenze. **Beide Wege gefahren:**

| Weg | Lauf | Ergebnis |
|---|---|---|
| Start erkannt | `make homebrew-nachweis` (der Beleg oben) | die Zeile des Wartens, dann `SIGTERM` nach einer Sekunde, drei `level=INFO`, Exit `0`; Dauer 14,4 s statt 17,5 s mit `sleep 5` |
| Grenze abgelaufen | Sonde: Kopie des gelieferten Skripts, in der die Schritte (a) bis (d) entfernt, die Grenze auf 2 Durchläufe gesetzt und der Lauf durch `bash -c 'trap "exit 0" TERM; sleep 30'` ersetzt ist — er schreibt keine Zeile und endet nach `SIGTERM` mit `0` | `homebrew-nachweis: die Grenze von 2 s ist abgelaufen, die Fehlerausgabe des Replay-Laufs ist leer geblieben`, dann `+ kill -TERM <Replay-Lauf>`, leere Fehlerausgabe, `--- Exit-Code des Replay-Laufs: 0`, `… (Exit-Code 7)`, `Skript endet mit Exit-Code 7`, Exit `7` — der Code des Laufs war `0`, das Skript endet dennoch mit `7`, genau wie Punkt 20 (3) es verlangt |

Die Sonde ist ein einmaliger Prüflauf, keine Gegenprobe und kein Artefakt im Repo (Punkt 20 (7): dieses Ziel hat keinen Fall; die Entscheidung des Nutzers vom 2026-10-21 schließt eine Gegenprobe aus). Sie lag in `/tmp` des Dev-Containers und ist gelöscht.

**F-2 (MEDIUM, behoben, Commit `74e19478`; Vertragstext des Architect in Punkt 20 (5)).** Der Kopf des Skripts sagt jetzt: Der Aufruf `--schritte` am Host gehört dem Betreiber, ein Agentenlauf fährt ihn nicht, sein Weg ist `make homebrew-nachweis` im Container — und ausdrücklich, dass das eine Regel an den Leser und kein Mechanismus ist: Der Wächter blockt ein Segment, dessen Kopf ein Paketmanager ist, liest den Inhalt einer Datei nicht und verhindert diesen Aufruf darum nicht. Damit sagt der Kopf keinen Schutz mehr zu, den der Wächter nicht leistet (`AGENTS.md` §3.11). Die Messung, die das zeigt, hat der Reviewer gefahren und der Architect in Punkt 20 (5) und §6 Randform 8 festgehalten; dieser Nachtrag wiederholt sie nicht.

**F-4 (LOW, behoben, Commit `47b1e489`).** Der Kopf des Handbuchs steht auf `Version: 0.10`. Die Reihe `git log -L3,3:docs/user/benutzerhandbuch.md` zählt je Doku-Zuwachs eine Stufe (0.1 bis 0.9); eine Regel für diese Zeile führt das Repo nirgends, deshalb setzt `0.10` die Zählung in der Lesart fort, die hier überall gilt (Stellen werden numerisch verglichen — `SPEC-067` Punkt 19 (9) führt dafür den Fall `p19-tap-semver-numerisch`: `0.10.0` ist höher als `0.9.0`). `1.0` wäre eine Aussage über Vollständigkeit, die niemand getroffen hat. `Stand: 21.10.2026` steht bereits auf dem Tag dieses Slice. **§11 *Änderungshistorie* bleibt unberührt:** Ihre Zeilen sind die Release-Notizen je veröffentlichter Software-Version (`SPEC-067` Punkt 3 (5), Punkt 9, Punkt 18); dieser Slice liefert kein neues Verhalten des Binaries `0.1.1`, sondern beschreibt einen Bezugsweg, der mit dem Release entstand.

**F-5 (LOW, behoben).** Der Satz über die Plattform des Docker-Servers trägt jetzt seine Befehle und Ausgaben (oben im Beleg): `docker version --format '{{.Server.Os}}/{{.Server.Arch}}'` gibt `linux/amd64`, `docker info --format '{{.OSType}}/{{.Architecture}}'` gibt `linux/x86_64`.

**F-7 (INFO, nichts zu ändern).** Am Endstand sind Plan und Diff deckungsgleich — das stellt der Report selbst fest. Die Reihenfolge der Commits (`e81a303c` Nutzer-Doku, dann `bcebc712` §1 und §6) bleibt stehen: Ein Commit, der §1 vorher nachzieht, hätte den Stand des Laufs behauptet, bevor der Beleg lag. `AGENTS.md` §3.9 verlangt die Nachführung „im selben Commit“ für Code, Gates, Tests, Spezifikation und Harness-Artefakte; ob Nutzer-Doku und ein Werkzeug ohne Gate darunter fallen, ist dem Wortlaut nach offen und eine Frage an den Architect, keine Änderung hier.

**F-9 (INFO, behoben).** `$(strip …)` steht jetzt unter den Entscheidungen des Implementers (Punkt 5 oben).

**§3 (Plan) nach dem Diff:** Die Tabelle nennt jede Datei, die dieser Slice geändert hat — `tools/harness/homebrew-nachweis.sh` und `harness/mk/homebrew-nachweis.mk` (neu), `docs/user/benutzerhandbuch.md` und `README.md` (update); `harness/README.md`, `harness/vertraege/SPEC-067-release.md` und `spec/spezifikation.md` lagen schon und sind unberührt. Keine Zeile fehlt, keine ist gegenstandslos geworden.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten** (Planner, 2026-10-21, gegen den gemergten Stand `origin/main` 28c5c1fe; die Sichtung vom 2026-10-03 fand ein Register mit nur seiner `README.md`). Das Register trägt heute **117 Einträge** unter `BEO-REPO/`. **Die Sub-Area ist `REPO` und umfasst das gesamte Repo; daraus folgt nicht, dass alle 117 Einträge diesen Slice treffen:** Das Kürzel sagt, wo ein Eintrag gezählt wird, nicht, dass sein Gegenstand in jedem Slice vorkommt.

**Maßstab der Auswahl — zielgerichtet, nicht erschöpfend.** Geprüft ist, was durch eines von zwei Sieben fällt:

1. **Zählerstand ≥ 2** — der Eintrag könnte mit diesem Slice die Schwelle erreichen. Das sind 51 Einträge.
2. **Gegenstand von diesem Slice berührt**, unabhängig vom Zähler: Release und Pipeline, Tap und Homebrew, Werkzeuge mit Netz oder fremdem Image, Handbuch und Ist-Zustand, Plattform-Grenzen, Belege und Angaben ohne Ausgabe.

**Nicht einzeln geprüft** sind die übrigen Einträge mit Zählerstand 1, deren Gegenstand außerhalb dieser Themen liegt (Protokollrand und Aufzeichnungsformat, Abdeckung und Laufzeit der Unit-Tests, Dev-Container und Span-Erfassung, CI-Range und Workflow-Pins). Grund: Sie erreichen mit diesem Slice keine Schwelle, und ihr Gegenstand kommt in seinen Dateien nicht vor. Tritt einer beim Lauf doch auf, ist das ein Auftreten wie jedes andere und gehört in die Closure.

**Sieb 1 — die 51 Einträge mit Zählerstand ≥ 2.** 27 tragen schon einen Ausgang und erreichen durch diesen Slice keine neue Schwelle (16 `verkörpert`, 8 `geplant` mit Kennung, 2 `gestrichen`, 1 `behoben`). Als **Regel** binden diesen Slice davon: `zusage-im-kommentar-weiter-als-pruefung` (89×, `AGENTS.md` §3.11 — Handbuch und Beleg sagen nur zu, was §7 zeigt), `plan-folgt-korrektur-nicht` (76×, §3.9), `negativtests-fehlen-bei-neuem-vertrag` (64×, §3.10 — hier durch Punkt 20 (7) ausdrücklich nicht einschlägig, §6 Randform 13), `spec-randform-erst-im-review-entschieden` (47×, §3.12 — der Vertragstext liegt vor dem Code), `folge-slice-adresse-nimmt-nicht-an` (13×, §3.13), `slice-waechst-durch-uebernahmen` (7×) und `schichtteilung-je-plan-verschieden` (5×, beide in den Nachzählungen unten), `groessenprobe-ohne-messbaren-stand` (5×, §3.14 — die Marke in §4 trägt ihre Metrik), `verhalten-nur-unter-linux-geprueft` (3×, §3.12 Plattform — §6 Randform 9 zitiert ihn).

Die übrigen 24 stehen auf `offen`. Zehn davon berührt dieser Slice:

| Eintrag (`BEO-REPO/…`) | Zähler | Warum er diesen Slice trifft |
|---|---|---|
| `gegenprobe-haengt-an-umgebung-des-hosts` | 9 | Der Nachweis läuft in zwei Umgebungen: im Container (Linux, Docker) und am Rechner des Betreibers (macOS, bash 3.2). §6 Randform 8 nennt die Grenze; die Klasse hat keinen Ausgang. |
| `gegenprobe-nicht-in-kette-regression-unbemerkt` | 4 | Der Slice legt ausdrücklich ein Werkzeug ohne Gegenprobe und ohne Kette an (§1 Abgrenzung, §6 Randform 13); eine Regression des Skripts bemerkt nur der Leser der Ausgabe. |
| `test-faellt-unter-cpu-last-aus` | 4 | Der Slice fährt `make gates`; `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` fiel viermal unter Last aus. Über der Schwelle, ohne Ausgang und ohne Adresse. |
| `nacharbeit-nach-review-ohne-erneuten-lauf` | 3 | Der Slice liefert ein ausführbares Skript, dessen Lauf Netz und den echten Tap braucht; Nacharbeit ohne erneuten Lauf liegt hier besonders nah. |
| `wiederanstoss-eines-roten-release-laufs-unbelegt` | 2 | Wiederanlauf des Jobs `homebrew` und Nachtrag von Hand gehören diesem Slice (§1 Bindung); nutzt er sie, ist das das dritte Auftreten. |
| `release-pipeline-pin-und-werkzeug-altern-ohne-melder` | 2 | Der Slice setzt einen neuen Pin (`homebrew/brew`), dessen Alterung kein Gate und kein Nachtlauf meldet (§6 Randform 3). |
| `folge-slice-liefert-handbuch-teil-nicht` | 2 | Der Slice trägt seinen Handbuch- und README-Teil selbst (Begründung in der letzten Nachzählung unten); bleibt der Teil aus, ist das das dritte Auftreten. |
| `historie-zeile-nennt-neuen-satz-der-spezifikation-nicht` | 2 | `LH-FA-19.a` hat mit dem Vertragstext einen neuen Absatz bekommen; die Änderungshistorie der Spezifikation muss ihn als neue Aussage führen, nicht als Berichtigung. |
| `abnahme-ohne-postgres-nicht-einzeln-brechbar` | 2 | Der Replay-Lauf des Nachweises läuft ohne PostgreSQL und ohne Fall (Punkt 20 (3), (7)); brechen kann ihn nur der Leser der Ausgabe. |
| `meldeschwelle-ohne-artefakt` | 2 | **Gegenbeleg:** §4 nennt Meldeschwelle (700) und Marke (900) im Plan — genau das Artefakt, dessen Fehlen der Eintrag beschreibt. |

Die vierzehn übrigen offenen Einträge mit ≥ 2 sind durchgegangen und treffen nicht, weil ihr Gegenstand in diesem Slice nicht vorkommt: `liste-gruener-mutanten-unvollstaendig` (kein Umbau, keine Mutanten), `zeilenende-crlf-im-gelesenen-wert`, `yaml-meldung-nennt-wert-und-zeile`, `plan-nennt-test-der-nicht-existiert` (der Plan nennt keinen), `postgres-versionsliste-altert`, `lastabhaengige-zaehlung`, `gegenprobe-startet-sich-ueber-einen-fall-selbst`, `gate-uebergeht-ablage-eintrag-still` (der Slice legt kein Gate und keine Ablage an), `stop-hook-positivliste-ohne-docs-developer`, `stop-hook-gegenprobe-prueft-nicht-die-umgebung`, `hostschluessel-nach-neubau-ohne-entscheidung`, `schnitt-laesst-haelfte-an-der-grenze` und `rueckfuehrung-ohne-verzeichniswechsel` (beide nur bei einer Rückführung, die §4 vorab benennt), `implementer-ueberschreibt-plan-abschnitte` (Klasse jedes Slice, kein Gegenstand dieses).

**Sieb 2 — Zählerstand 1, Gegenstand berührt.** Neun Einträge, alle `offen`:

| Eintrag (`BEO-REPO/…`) | Zähler | Warum er diesen Slice trifft |
|---|---|---|
| `vorbedingung-geheimnis-nur-existenz-geprueft` | 1 | Das **Recht** des Tap-Geheimnisses prüft die Checkliste nicht, nur seine Existenz; §6 führt genau dieses Risiko (Schreibrecht ohne Eingriff des Betreibers weiter unbelegt). |
| `dod-beleg-vom-belegverfahren-nicht-erzeugbar` | 1 | Die DoD verlangt Belege von Plattformen, die der Implementer nicht erreicht. Die Kombinations-Lesart (§6 Randform 9, 10) ist das Gegenmittel: kein Beleg, keine Zusage, dafür eine Grenze — der Slice bleibt schließbar. |
| `vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen` | 1 | Punkt 20 (1) und (2) legen Aufruf und Umgebung von Homebrew im fremden Image fest, ohne dass ein Lauf sie gemessen hat; §6 führt das als Risiko mit Rückgabe an den Architect. |
| `vertrag-nennt-exit-code-des-skripts-nicht-von-make` | 1 | Punkt 20 (2) nennt den Exit-Code des **Skripts**; über `make` sieht der Aufrufer 2. Beim Schreiben des Fragments zu prüfen, damit der Beleg nicht den falschen Code nennt. |
| `docker-run-ohne-rm-ausserhalb-der-laeufer` | 1 | Der Slice legt ein neues `docker run` an; Randform 5 verlangt `--rm` samt anonymen Volumes. |
| `beleg-ueber-den-bestand-ohne-probe` | 1 | Der gemessene Stand des Jobs `homebrew` in §6 ist eine Aussage über den Bestand; die Probe dazu ist der eigene Lauf dieses Slice. |
| `empfehlung-als-nutzerentscheidung-bezeichnet` | 1 | **Gegenbeleg:** Die Angabe des Nutzers zum macOS-Lauf steht in §7 als Angabe der Koordination — weder als Entscheidung noch als Beleg. |
| `anforderung-ohne-zeile-in-der-abdeckungstabelle` | 1 | `LH-FA-19` trägt keine Test-Deklaration; der Nachweis ist kein Test und erzeugt keine, der Abgleich meldet dazu nichts. |
| `werkzeug-festlegung-ausserhalb-technik-stratum` | 1 | **Gegenbeleg:** Die Festlegung des Werkzeugs steht als Punkt 20 in `harness/vertraege/SPEC-067-release.md` (`MR-002`), nicht in einer ADR und nicht im Skript. |

**Kriterium Evidenz-/Diskrepanz-Risiko** (Pflichtkriterium 3 dieses Abschnitts, für die Sub-Area `REPO`): **hoch für den Lauf, niedrig für den Text.** Der Vertragstext liegt vor dem Code (`spec-randform-erst-im-review-entschieden` verkörpert), die Größe ist nachgezählt, die Plattform-Grenze ist entschieden statt offen (`verhalten-nur-unter-linux-geprueft` verkörpert). Das Risiko liegt im **Lauf gegen fremde Dienste**, und dafür stehen vier Einträge über der Schwelle ohne Ausgang — `gegenprobe-haengt-an-umgebung-des-hosts` (9×), `test-faellt-unter-cpu-last-aus` (4×), `gegenprobe-nicht-in-kette-regression-unbemerkt` (4×), `nacharbeit-nach-review-ohne-erneuten-lauf` (3×) — und zwei unter ihr, deren Gegenstand der Slice unmittelbar berührt: `vorbedingung-geheimnis-nur-existenz-geprueft` (1×, die Rechte des Tap-Tokens) und `vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen` (1×, die ungemessenen Annahmen über das Image).

**Schwelle.** Die Sichtung ist kein Auftreten; kein Eintrag erreicht durch sie drei. **Fünf Einträge stehen bei 2, und ihr Gegenstand liegt in diesem Slice** — `wiederanstoss-eines-roten-release-laufs-unbelegt`, `release-pipeline-pin-und-werkzeug-altern-ohne-melder`, `folge-slice-liefert-handbuch-teil-nicht`, `historie-zeile-nennt-neuen-satz-der-spezifikation-nicht`, `abnahme-ohne-postgres-nicht-einzeln-brechbar`. Tritt einer von ihnen beim Lauf auf, erreicht er mit der Closure dieses Slice 3× und braucht dann einen eigenen Folge-Slice mit Kennung (Baseline-Regelwerk `modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung); die Closure trägt das ein, und der Planner legt den Slice erst auf Anweisung an. Die vier Einträge **über** der Schwelle brauchen keinen Folge-Slice aus diesem Slice, sondern einen Ausgang, den der Nutzer zuweist; der Lese-Schritt der Closure von [welle-v1-nachtrag](../welle-v1-nachtrag.md) sichtet sie.

Die Sichtung ändert Liefer-Punkte und Schichten nicht: Sie ist Pflichtinhalt dieses Abschnitts, keine Sendung (`AGENTS.md` §3.13).

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (Tap mit Formel, Dokumentation). Nachgezählt am 2026-10-20 mit der Sendung aus `slice-v1-abschluss-homebrew` (Architect, `AGENTS.md` §3.13): **Liefer-Punkte: 2** (1: Formel im echten Tap; 2: Abnahmeszenario 11 mit Handbuch- und README-Teil); Schichten bleiben zwei (der Tap als Ergebnis des Jobs oder einer Handlung des Betreibers, Dokumentation). Nachgezählt am 2026-10-20 mit der Bindung aus `slice-v1-abschluss-binary-namen`: Der Nachtrag für `v0.1.0` entfällt, die Formel entsteht durch den Job; **Liefer-Punkte bleiben 2, Schichten zwei**; der Slice liegt unter der Grenze.

Nachgezählt am 2026-10-20 mit der Bindung aus `slice-v1-abschluss-release-v0-1-1` (Planner, `AGENTS.md` §3.13): Der Start-Trigger nennt den Slice, der Tag kommt von ihm; die Abgrenzung der Handbuch-Abschnitte (Binary und Image beim Release-Slice, Homebrew hier) verkleinert den Teil dieses Slice; **Liefer-Punkte bleiben 2, Schichten zwei**; der Slice liegt unter der Grenze.

Nachgezählt am 2026-10-20 mit der Sendung aus `slice-v1-abschluss-homebrew` (Wiederanlauf des Jobs, Planner, Closure jenes Slice, `AGENTS.md` §3.13): ein Risiko in §6, kein Liefer-Punkt und keine Schicht; **Liefer-Punkte bleiben 2, Schichten zwei** (der Satz in `releasing.md` ist Wartungs-Doku, die der Slice schon trägt); der Slice liegt unter der Grenze.

Nachgezählt am 2026-10-21 mit dem Risiko-Ausgang aus `slice-v1-abschluss-release-v0-1-1` (Planner, Closure des Gebers, `AGENTS.md` §3.13): Eingetragen sind der gemessene Stand des Jobs `homebrew` und die Liste dessen, was weiter unbelegt ist — beides steht in §6 und gehört zum ersten Liefer-Punkt (Formel im echten Tap) und zum zweiten (Abnahmeszenario 11); es kommt kein Liefer-Punkt und keine Datei hinzu, und die Berichtigung in §1 und §4 (den Tag setzte die Koordination, nicht der Nutzer) ändert nur die Beschreibung des Trägers. **Liefer-Punkte bleiben 2, Schichten zwei** (Tap mit Formel, Nutzer-Doku); der Slice liegt unter der Grenze.

Nachgezählt am 2026-10-21 vom Architect mit der Entscheidung der Randformen (`AGENTS.md` §3.13): **Liefer-Punkte 3** — (1) die Formel des stabilen Tags im echten Tap; (2) das Werkzeug `make homebrew-nachweis` (Fragment, Skript, Zeile in `harness/README.md`) mit seinem Lauf unter Linux; (3) Abnahmeszenario 11 je Kombination belegt oder begrenzt, mit dem Homebrew-Abschnitt in §2 *Installation* und der Zeile in `README.md`. **Schichten 2** — Harness-Werkzeug (`tools/harness/`, `harness/mk/`, `harness/README.md`) und Nutzer-Doku (`docs/user/benutzerhandbuch.md`, `README.md`); nach der Teilung von `slice-harness-d-check-v0-85` §8 zählen Verträge (`harness/vertraege/`, `spec/`) und Planung nicht. **Korrektur der früheren Zählungen dieses Plans** (sie führten „Tap mit Formel“ als Schicht, zuletzt die Zählung darüber): Der Tap ist ein fremdes Repository, das der Job `homebrew` schreibt; dieser Slice ändert dort nichts, er belegt — eine Schicht ist eine berührte Schicht dieses Repos. Damit liegt der Slice unter beiden Grenzen, und der Handbuch-Teil bleibt hier: ein Doku-Folge-Slice `slice-erster-release-homebrew-nachweis-doku` (`AGENTS.md` §3.11, §3.13) ist nicht nötig.

Nachgezählt am 2026-10-21 beim Nachziehen des Plans auf den Stand des Architect und beim Claim (Planner, `AGENTS.md` §3.13): Eingetragen sind die beiden Angaben der Koordination in §7 (macOS-Lauf des Nutzers, fehlende Option `--version`) — sie sind Eingaben für den Implementer, kein Liefer-Punkt und keine Datei; dazu die Nachführung von DoD-Punkt 1 auf die Lesart je Kombination (`LH-FA-19.a`, `SPEC-067` Punkt 20 (6)), die den Punkt nicht vergrößert, sondern seine Zusage an den Vertragstext bindet. **Liefer-Punkte bleiben 3, Schichten bleiben 2** (Harness-Werkzeug, Nutzer-Doku); der Slice liegt unter beiden Grenzen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (`harness/conventions.md` deklariert eine Sub-Area für das gesamte Repo, Modus Greenfield).
