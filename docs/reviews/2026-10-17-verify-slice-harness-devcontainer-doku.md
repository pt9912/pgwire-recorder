# Verifikationsbericht: slice-harness-devcontainer-doku — 2026-10-17

**Art:** Verifikation (Modul 11) gegen DoD, Plan und Belege. Fragestellung: Bauen wir es richtig?
**Gegenstand:** Stand `affbd5a2` (Dokument `531f357e`, `d725a22c`; Review `c5576654`; Plan `affbd5a2`).
**Rolle:** Verifier, frischer Kontext. Nichts repariert, nichts am Container verändert.

## Verdikt

**Bestanden mit Vorbehalt.** Das Dokument stimmt mit dem Versuchs-Container und dem Plan überein;
die F-1-Korrektur hält. Offen sind zwei Punkte an den Planner (Form der Risiko-Ausgänge, Beleg der
Gates nach der Nacharbeit) und ein Hinweis zur Meldeschwelle. Kein HIGH.

## Eigene Läufe

- `diff` der Codeblöcke 3 bis 5 (`Dockerfile`, `start.sh`, `compose.yaml`) gegen
  `/Development/pgr-dev/.devcontainer/` und `compose.yaml`: **kein Unterschied**.
- Lesende Proben: `docker port pgr-dev` = `127.0.0.1:2222`; `ss -ltn` nur `127.0.0.1:2222`;
  Mounts genau vier (`/host-repo` rw=false, `/workspaces`, `authorized_key.pub` rw=false,
  `/var/lib/docker`), `Privileged=true`; per SSH `sudo -n -l` = `(root) NOPASSWD: ALL`,
  `mount` zeigt `/host-repo` als `ro`, `id` = uid 1000, Gruppe `docker`.
- `git diff --numstat bda6efef..HEAD -- docs/developer` = 407 / 0, Summe **407**. Basis stimmt:
  `bda6efef` ist der Elternteil von `0657b0e5` (next -> in-progress).
- **Nicht gefahren:** `make docs-check`, `make kopf-check`, `make gates`. Es läuft ein voller
  `make gates` (`pgrep`: PID 2209712), also kein `docs-check` (Auftrag).

## DoD

| Punkt | Befund |
|---|---|
| Dokument deckt §1 ab, Schritte nachgeführt oder „nicht geprüft“, Dateiinhalte gleich | erfüllt: alle Punkte aus §1 vorhanden (§§2 bis 12), dazu §13; Dateiinhalte per `diff` gleich; Probentabelle in §7 deckt die Aussagen |
| `kopf-check`, `docs-check`, `make gates` grün | **nicht bestätigt**: Belege in §7 (`docs-check` 0 mit 966 Dateien, `kopf-check` 0) stammen aus `531f357e`, vor der Nacharbeit `d725a22c`; nach ihr ist kein Lauf belegt. Voller `make gates` läuft (Stand heute); sein Ergebnis steht aus |
| Review, Report liegt vor | erfüllt (`docs/reviews/2026-10-17-review-slice-harness-devcontainer-doku.md`, Ausgang je Finding in §7) |
| Closure-Notiz, Lerneintrag | offen (Planner, Closure; §7 trägt noch Platzhalter *Wird bei der Closure gefüllt*) |
| Register | offen (Planner) |
| Risiko-Ausgänge | siehe V-2 |
| Paarungen | offen (Planner) |

Alle DoD-Häkchen stehen noch leer; das ist vor der Closure zu erwarten.

## Prüfung F-1 (Sicherheitsaussagen)

- `grep schreibgesch` : **keine Treffer**. Einziger Treffer für „nicht schreiben“ (Z. 307): „Ein Prozess ohne
  `sudo` kann dort nicht schreiben“ — eng gefasst, belegt durch `touch` in §7 (Read-only file system).
- §1 (keine Sandbox, `sudo` ohne Passwort), §3.2 (Schlüssel ohne Passphrase, Folge), §9 (kein Remote, Mount `ro`,
  Commits prüfen), §11 (privilegiert, `sudo`, „nicht ausgeführt, nur gefolgert“) stimmen untereinander und mit
  meiner Probe überein. „nur Loopback“ und „nur Schlüssel“ sind belegt (`docker port`/`ss`; drei abgewiesene
  Anmeldungen laut §7, von mir nicht wiederholt, da das Anmeldungen mit fremdem Schlüssel verlangte).
- „genau vier Mounts“ (§7), „Docker-Socket nicht gemountet“ (§11): durch `docker inspect` bestätigt.
- Die Quantoren „nie/immer/stets/jede/alle“ kommen im Dokument nur in unkritischen Zusammenhängen vor
  (Z. 338 „wie jeder andere“, Z. 344 „jeder im eigenen Verzeichnis“ = Beschreibung der Probe).
- Aussagen „kein Remote zeigt auf ihn“ ist auf das Host-Repo am Prüftag eingegrenzt; ich habe `git remote -v`
  nicht wiederholt (Beleg liegt im Plan).

## „Nicht geprüft“ (§13)

§13 des Dokuments nennt zwölf Punkte und deckt sich mit der Liste in Plan §7 (acht ursprünglich, vier seit der
Nacharbeit: Hostschlüssel-Frage, Fetch ins echte Host-Repo, Schreiben mit `sudo`, `docker compose up -d` allein).
Jeder Verweis (§4, §5, §6, §9, §11, §12) auf §13 hat einen Gegenpunkt im Text. **Stimmt.**
Gegen die Behauptung „die Zeiten des Versuchs-Containers nicht wiederholt“ steht §8 offen mit Quelle; in §13
fehlt dieser Punkt als eigene Zeile, er steht aber in §8 selbst (siehe V-4).

## Plan folgt Diff (§3.9)

§1 (Ziel, Abgrenzung), §3 (Nacharbeit, Quellen, Lücken), §6 (Risiken) und §7 folgen dem Dokument. Abweichungen:

- **V-4 (INFO):** Plan §3 sagt, das Dokument nenne keine Kennungen von Lastenheft, Spezifikation, ADRs. Es
  enthält in den Codeblöcken `LH-FA-DEV-004`, `LH-FA-ADD-*`, `LH-SA-FILE-002` (Z. 94, 153, 154). Das sind die
  Kennungen von `u-boot` im wörtlich übernommenen Dateiinhalt, nicht des Repos, haben aber die Form `LH-…`.
  Ob `.d-check.yml` (`ids`, `link-policy: always`) sie als Kennung dieses Repos liest, konnte ich ohne
  `docs-check` nicht klären; der Beleg `docs-check` 0 stammt aus `531f357e`, wo sie schon standen.
  Der Plan sollte den Satz auf „keine Kennungen dieses Repos außer den wörtlich übernommenen Dateiinhalten
  von `u-boot`“ nachziehen.
- Plan §7 „Läufe“ nennt noch die Größe 385 (Stand `531f357e`) neben 407 weiter unten; der Verlauf ist
  nachvollziehbar, aber die ältere Zahl steht ohne den Hinweis „vor der Nacharbeit“ in derselben Zeile wie der Lauf.
  Geringfügig.

## Größe

407 (Einfügungen 407, Löschungen 0), Basis `bda6efef`. Marke 500 nicht erreicht, **Meldeschwelle 400 überschritten**.
Plan §4 nennt die Schwelle, definiert aber keine Folge; Plan §7 (Nacharbeit) vermerkt „über der Meldeschwelle 400,
unter der Marke 500“ als Zahl, ohne Meldung oder Entscheidung. Das ist hiermit die Meldung: **Planner nimmt sie in
die Closure auf** (Lerneintrag: Schätzung 150 bis 300 lag 36 bis 170 % unter der Messung; §3.14 verlangt Schätzungen in
dieser Metrik, hier fehlte der Faktor für Prosa mit Dateiinhalten und Probenbeschreibung). Rückführung nach §4 trat
nicht ein.

## Stop-Hook / Positivliste (Risiko 1)

Korrekt beschrieben: `docs/developer/` steht nicht auf der Positivliste (`docs/plan/planning/`, `docs/reviews/`,
`docs/user/`, Wurzel-`README.md`), der Hook verlangt einen vollen `make gates` nach dem Stand dieses Arbeitsbaums.
Ein solcher läuft gerade (PID 2209712); sein Ergebnis und ein Stand der Dateiliste, der den Commit `d725a22c` und
alle späteren umfasst, sind vor dem Push erforderlich. Gelangen danach weitere Commits unter `docs/developer/`,
ist der Lauf zu wiederholen.

## Befunde

| Id | Schwere | Befund |
|---|---|---|
| V-1 | MEDIUM | Der Beleg von `make docs-check`/`kopf-check` ist älter als die Nacharbeit `d725a22c`. Nach ihr ist kein Lauf in §7 belegt; ich durfte keinen fahren. Die DoD verlangt grün auf dem Endstand. Erst der laufende volle `make gates` (oder ein eigener Lauf nach Ende) belegt es. |
| V-2 | MEDIUM | Risiko-Ausgänge in §6 sind nicht eine der drei geschlossenen Formen. Risiko 1: „weiter offen: → Entscheidung des Nutzers“ — *weiter offen* heißt Eintrag im Beobachtungs-Register (Modul 5 §Offene Risiken), nicht eine Entscheidung. Risiko 2: „weiter offen … falls …; sonst entfallen“ ist bedingt, also nicht ein Ausgang. Beides muss der Planner bei der Closure auf **einen** Ausgang je Risiko festlegen. |
| V-3 | LOW | Meldeschwelle 400 überschritten (407), ohne benannte Folge im Plan (siehe Abschnitt Größe). |
| V-4 | INFO | Plan §3 „keine Kennungen“ widerspricht den `LH-…`-Kennungen von `u-boot` in den Codeblöcken (siehe oben). |
| V-5 | INFO | Die in §8 genannten Zeiten des Versuchs-Containers (5 s, 201 s) stammen aus Skripten des Nutzers, nicht wiederholt; offengelegt in §8, aber nicht in der §13-Liste. Kein Verstoß gegen §3.11. |

## Pflichten des Planners bei Closure

1. **§7 Closure-Notiz** füllen (Platzhalter ersetzen); Lerneintrag in einer der drei Formen
   (geschärfte Regel · neuer Sensor · benannte Spec-Lücke). Kandidaten: Schätzung in der Metrik von §3.14
   unterschätzt auch Prosa-Dokumente (407 gegen 150 bis 300); `docs/developer/` ohne Positivliste und ohne
   Sensor, der das Dokument gegen den Container hält.
2. **Risiko-Ausgänge §6** je genau einer (V-2): R1 *weiter offen* → Register oder *eingetreten* mit Kennung;
   R2 *entfallen* mit Begründung oder *weiter offen* → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`;
   R3 *entfallen* (Dokument markiert „nicht geprüft“) stimmt; R4 ist kein Risiko.
3. **Register** unter `docs/plan/planning/observations/` fortschreiben: Zähler der Beobachtungen, die dieser Slice
   berührte (Finding-Klassen des Reviews: Aussage über die Prüfung hinaus [F-1, F-4 → `zusage-im-kommentar-weiter-als-pruefung`];
   F-2/F-3 Quelle und Lücken; F-7 Plan/Kopf) mit dem Pfad `BEO-REPO/<slug>`; sonst „keine Beobachtung“ in §7. Der
   Zähler wird aus den Dateien gelesen (§8 des Plans), nicht aus dem Gedächtnis.
4. **Finding-Klassen zählen:** die wiederkehrende Klasse aus der Summary-Zeile des Review-Reports (1 MEDIUM, 4 LOW,
   2 INFO) geht in den Closure-Eintrag.
5. **Paarungen** (Anker · Folge-Slice · Register): jeder neue Anker `seit slice-harness-devcontainer-doku` hat
   Gegenstück; eine Entscheidung zur Positivliste (Nutzer) braucht eine Adresse (Folge-Slice oder Register-Zeile), die
   vorher annimmt (`AGENTS.md` §3.13).
6. **Reihenfolge:** erst Inhalt (DoD-Häkchen, Notiz, Verifier-Bericht), dann reiner `git mv` nach `done/` (`AGENTS.md` §3.3);
   `make gates` auf dem Endstand grün, bevor gepusht wird (Hook).
7. **Plan nachziehen (§3.9):** Satz zu Kennungen (V-4); Zeile „Läufe“ in §7 mit Stand der Zahl.

## SHA

Geprüfter Stand: `affbd5a2`. Dieser Bericht wird mit `git commit --only -- docs/reviews/2026-10-17-verify-slice-harness-devcontainer-doku.md`
eingecheckt; die SHA des Berichts-Commits nennt die Übergabe.
