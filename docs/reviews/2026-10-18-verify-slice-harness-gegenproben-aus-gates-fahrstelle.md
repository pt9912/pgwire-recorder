# Verifikation slice-harness-gegenproben-aus-gates-fahrstelle (Teil A: Hook-Skript)

Stand: Zweig `rev-fahr`, Implementer-Stand `d5e15f14`, Basis der Metrik `0f066c04`, Vergleichsbasis des Plans `dff05377`.
Maßstab: [Plan](../plan/planning/done/slice-harness-gegenproben-aus-gates-fahrstelle.md) §2 und §7,
[SPEC-061](../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md), [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md).
Host-Läufe als Benutzer `db` (uid 1000, kein root), `LANG=de_DE.UTF-8`, Host-Zeitzone CEST.

## Ergebnis je Punkt

| # | Punkt | Ergebnis |
|---|---|---|
| 1 | `make gates-schnell-gegenprobe` unter `LANG=de_DE.UTF-8` | grün, 1 min 03 s |
| 2a | Gegenprobe des Hooks, ganz, am Host (`make gegenproben-nach-einlieferung-gegenprobe`, de_DE) | **ROT** (4 min 19 s): Mutation `m-p9-ende-lokal` färbt zusätzlich `p9-beginn-nach-sperre` (Befund F-1) |
| 2b | dieselbe Gegenprobe mit `TZ=UTC` für Hook und Attrappe (Kopie in Temp, zwei Zeilen geändert) | grün, 74 Fälle, 113 Mutationen, 4 min 17 s |
| 2c | 12 Mutanten einzeln per `GP_MUT`/`GP_ENTDECKE` (Punkte 2 bis 13) | alle rot mit genau der erwarteten Menge, grün gemeldet, 27 s |
| 3 | echter Lauf (Bare-Repo, Hook per `install -m 755`, Klon mit Makefile) | bestanden, siehe unten |
| 4 | Zweignamen mit Shell-Metazeichen | sicher, siehe unten |
| 5 | nicht als root | bestanden mit Anmerkung, siehe unten |
| 6 | `docs-check`, `kopf-check`, `gates-schnell`, Metrik | `docs-check` 0 Befunde, `kopf-check` Exit 0, `gates-schnell` Exit 0, Metrik 701 |

### Zu 2c

Gefahren: `m-p2-absolut` (`p2-relativ`), `m-p3-rueckwaerts`, `m-p4-41-stellen` (`p4-hash-41`), `m-p5-erbt-stdin` (`p5-fd-0`),
`m-p6-ohne-sperre` (`p6-nacheinander p9-beginn-nach-sperre p9-sigterm-wartet`), `m-p7-altes-bleibt` (`p7-altes-weg`),
`m-p8-grenze-121` (`p8-lang`), `m-p9-rot-rueckwaerts` (`p9-rot`), `m-p10-host` (`p10-argv`), `m-p11-makeflags-bleibt`
(`p11-makeflags`), `m-p12-bleibt-nach-rot` (`p12-entfernt-rot`), `m-p13-ohne-ablage-probe` (`p13-ablage-vor-zeile p13-ablage-weg`).
Die Gegenprobe verlangt je Mutation genau die benannte Menge und meldete grün. Abweichung zur Tabelle in §7 (F-3):
`m-p3-rueckwaerts` färbt `p12-stand p3-mehrere p9-gruen p9-rot`, die Tabelle nennt nur `p3-mehrere`; die Tabelle ist gekürzt, das
Skript trägt die genaue Menge.

### Zu 3 (echter Lauf, Temp unter dem Scratchpad)

- Push von `feature/x` kehrt nach 0,033 s zurück; der Hook meldet auf stderr `gestartet: feature/x <hash>`; Hook hält die Leitung nicht.
- Marke `gegenproben-feature%2Fx-<hash>.laeuft` und `.log` stehen sofort da; nach dem Ende liegt die Ergebnisdatei
  `gegenproben-feature%2Fx-<hash>` mit `zweig=feature/x`, `hash`, `exit=0`, `ergebnis=gruen`, `beginn`, `ende` (UTC, `Z`), `rot=`;
  `.laeuft` ist weg, Arbeitsbaum weg.
- Derselbe Hash unter anderem Namen (`feature/y`) fährt; derselbe Hash unter demselben Namen nach Löschung und erneutem Push fährt neu
  (altes Ergebnis beim Aufruf entfernt, `.laeuft` entsteht, neues Ergebnis danach). `+feature/x` ohne Änderung löst nichts aus
  (`Everything up-to-date`, git ruft den Hook nicht).
- Löschen des Zweigs (`git push o :feature/x`): kein Lauf, keine Meldung, vorhandene Ergebnisdateien und Logs bleiben.
- Rotes Makefile (`gegenproben: lint: rot (Exit 2)`, Exit 2): `exit=2`, `ergebnis=rot`, `rot=lint`.

### Zu 4 (Zweignamen)

Git lehnt Leerzeichen, Tabulator, Newline, `\`, `?`, `*`, `[`, `~`, `^`, `:` im Ref ab; damit entfallen Newline und Leerzeichen. Erlaubt und
gepusht: `` ` ``, `$(…)` und `;` (ohne Leerzeichen über `${IFS}`), `'`, `"`, `&`, `|`, `<>`, `{}`, `()`, `$`, `-rf`, `ü/ö`. Kein Befehl lief
(`/tmp/PWNED*` fehlt), der Zweigname steht nur als Argument des Hintergrundlaufs und im `printf`-Argument der Ergebnisdatei; die Dateinamen
sind kodiert (`%60`, `%24`, `%3B`, `%C3%BC`), der Name `-rf` steht hinter `gegenproben-`. Nicht belegt: ein Zweig mit Newline, den git nicht zulässt.

### Zu 5 (root)

Der Gegenproben-Container läuft als `vscode` (uid 1000, `docker exec pgr-dev2 id`); der Hook läuft als der Benutzer, der pusht.
Die Gegenprobe endet als root bewusst mit Exit 2 und Meldung (Zeile 16), nicht mit einem falschen Grün. Als root nicht gefahren (kein root am
Host). Offen bleibt die Rückgabe (6) des Implementers: Bedingung in `SPEC-061` oder anderer Aufbau, Entscheidung des Architect.

### Zu 6 (Metrik und Plan-vs-Code)

`git diff --numstat 0f066c04..HEAD -- tools harness/mk harness/README.md`: 116 + 572 + 11 + 2 = 701 Einfügungen, 0 Löschungen, Summe **701**
(wie §7). Meldeschwelle 700 um 1 überschritten, Marke 900 nicht; keine Rückführung. Seit `dff05377` ändert der Plan nur §7 (60 Zeilen);
Code: zwei Dateien in `tools/harness`, ein Fragment, zwei Zeilen in `harness/README.md`, nichts darüber hinaus. `harness/README.md` führt die
Zeile der Gegenprobe (Sensors, Eintrag in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`) und die Zeile des Skripts (Werkzeuge, kein Gate).

## DoD-Konformität

- Liefer-Punkt 1 (Hook-Skript): bestätigt durch den echten Lauf (Punkte 1 bis 9, 11 bis 13 am Verhalten; Punkt 10 durch `p10-argv`).
- Liefer-Punkt 2 (Gegenprobe): bestätigt in der UTC-Umgebung (2b), in der Host-Zeitzone rot (F-1).
- `make gates` grün: vom Implementer mit 4 min 19 s belegt, hier nicht wiederholt; `gates-schnell` grün, `docs-check` und `kopf-check` grün.
- Doku-Update `harness/README.md`: zwei Zeilen vorhanden. Die Zeile der Gegenprobe nennt die Werkzeuge nicht vollständig (`pgrep`, `timeout`
  stehen im Kommentar des Fragments, nicht in der Zeile; F-4, gering).
- Register, Review, Closure-Notiz, Risiko-Ausgänge, Paarungen: noch offen (Closure-Schritte, nicht Teil dieser Verifikation).

## Befunde

- **F-1 (mittel): Die Gegenprobe hängt an der Zeitzone des Hosts.** `haken()` und die Zeile 265 starten den Hook mit `env -i` ohne `TZ`; er liest
  dann `/etc/localtime`. Auf einem Host mit anderer Zeit als UTC färbt `m-p9-ende-lokal` zusätzlich `p9-beginn-nach-sperre` (lokales `ende`
  liegt hinter dem UTC-`beginn` des zweiten Laufs, `g6`), die Gegenprobe meldet ROT mit der Meldung „erwartet rot: p9-gruen; rot sind:
  p9-beginn-nach-sperre p9-gruen“. Im Container (UTC) grün, deshalb sieht es §7 nicht. Auswirkung: `make gegenproben-nach-einlieferung-gegenprobe`
  ist am Host unter CEST rot, `SPEC-060` Punkt 2 hält sie zwar am Container, aber die Zusage „Docker braucht sie nicht“ und das Review, das
  sie „gezielt selbst“ fährt, laufen am Host. Behebung (Implementer): `TZ=UTC` im `env -i` von `haken` und Zeile 265 setzen (belegt: damit
  grün, 113 Mutationen), oder die erwartete Menge um den Fall ergänzen. Hängt an keiner Spezifikationsstelle, kein Rückgabefall an den Architect.
- **F-2 (gering): `p13-ablage-weg` verlangt einen Benutzer ohne root.** Als root Exit 2, also `make gegenproben` rot in einem Root-Container.
  Der Gegenproben-Container läuft als `vscode`; Rückgabe (6) des Implementers an den Architect bleibt offen.
- **F-3 (gering): Tabelle in §7 kürzt rote Mengen** (`m-p3-rueckwaerts`); „alle 113 Mutationen mit ihrer genauen roten Menge stehen als
  `mut`-Zeilen“ ist richtig, die Tabelle ist es nicht in jeder Zeile.
- **F-4 (gering): Werkzeugliste der README-Zeile** nennt `pgrep` und `timeout` nicht, die DoD nennt bash, `git`, `flock`, `setsid`.
- Die Rückgaben (1) bis (5) des Implementers sind Sache des Architect; keine habe ich entschieden. Der echte Lauf bestätigt (3): letzte Zeile ohne
  Zeilenumbruch wird verarbeitet.

Bewertung: Hook-Verhalten bestätigt (echter Lauf, Metazeichen, Randformen), Gegenprobe in UTC bestätigt; am Host mit Nicht-UTC-Zeitzone rot (F-1).
Der Slice ist nach Behebung von F-1 verifizierbar.

## Nachverifikation (Stand 0e2545ff)

Host-Läufe als `db` (kein root), Host-Zeitzone CEST, `LANG=de_DE.UTF-8`; Container nicht gefahren. Host-Last während der Läufe: Load 7 bis 9
(fremde Prozesse, u. a. `suricata` mit 100 % CPU).

| # | Punkt | Ergebnis |
|---|---|---|
| 1 | Gegenprobe des Hooks ganz am Host, zweimal | **ROT beide Male**, je 24 min 09 s (nicht 5 min); rot ist in 41 bzw. 55 der 131 Mutationen genau ein zusätzlicher Fall, `p9-sigterm-wartet-waise`; sonst stimmt jede rote Menge; der Fall der Zeitzone ist nicht beteiligt |
| 2 | 18 Mutanten einzeln per `GP_MUT`/`GP_ENTDECKE` | alle grün gemeldet, rote Menge genau wie in §7: `m-p7-alt-ungleich-null-entfaellt` (`p7-force`), `m-p8-ohne-feste-locale` (`p8-locale`), `m-p8-lang-ohne-meldung` (`p8-lang p8-lang-mehrbyte`), `m-p9-sigterm-ohne-trap` (`p9-sigterm-enkel p9-sigterm-lauf p9-sigterm-wartet-waise`), `m-p9-sigterm-wartet-unbehandelt` (`p9-sigterm-wartet p9-sigterm-wartet-waise`), `m-p9-beginn-lokal` und `m-p9-ende-lokal` (je `p9-zeitzone-ost p9-zeitzone-west`), `m-p11-zwei-namen` (`p11-git-weitere`), `m-p11-nur-git-dir` (`p11-git-quarantaene p11-git-weitere`), `m-p11-gnumakeflags-bleibt` (`p11-gnumakeflags`), `m-p13-werkzeug-umgekehrt` (`p13-werkzeug-erstes`), `m-p13-mktemp-beendet` (`p13-mktemp`), `m-p13-signal-ohne-ergebnis` (`p13-make-signal`), `m-p13-127-ohne-ergebnis` (`p13-make-127`), `m-p4-leerzeile-meldet` (`p4-leerzeichen-zeile p4-leerzeile`), `m-p4-vier-felder-ref` (`p4-vier-felder`), `m-p2-fehlt-nimmt-ablage` (`p13-marke-vor-werkzeug p2-fehlt`), `m-p2-leer-nimmt-ablage` (`p2-leer`) |
| 3 | echter Hook-Lauf (Bare-Repo, `install -m 755`, Marke, Ablage, Klon mit Makefile) | bestanden |
| 4 | `docs-check`, `kopf-check`, `gates-schnell` | 0 Befunde, Exit 0, Exit 0 |
| 5 | Metrik | `git diff --numstat 0f066c04..HEAD -- tools harness/mk harness/README.md`: 118 + 663 + 11 + 2 = **794** Einfügungen, 0 Löschungen (wie §7); Plan seit `dff05377`: zwei Hunks, §6 (+2) und §7 (+150), sonst nichts |

### Zu 3 (echter Lauf)

- `TZ=Asia/Kolkata` im Push: Push kehrt nach 0,033 s zurück; Ergebnisdatei `beginn=2026-10-18T14:12:03Z`, `ende=…14:12:06Z` (UTC, Host 16:12 CEST).
- Leere Zeile (und `\n` vor und nach der Zeile) auf stdin des Hooks: Zeile verarbeitet, Exit 0, Ergebnis angelegt.
- Force-Push (`-f` auf `feature/x`, anderer Hash): neuer Lauf, drei Ergebnisdateien nebeneinander, Zweig mit dem neuen Hash gefahren.
- 60 mal `ü`: abgelehnt mit „länger als 120 Zeichen“ (kodiert 360 Zeichen), Hook Exit 0, Push gelungen, kein Lauf; das ist `p8-lang-mehrbyte`.
- `GNUMAKEFLAGS=-i` im Pusher und rotes Makefile (Exit 2): `exit=2`, `ergebnis=rot`, `rot=lint`; mit grünem Makefile `exit=0`.
- mktemp-Fehler: Sperre von außen gehalten, zwei Zweige in einem Push, Ablage nach dem Anlegen der Sperre `a-w`, dann wieder schreibbar: Log des ersten
  „arbeitsbaum: Verzeichnis nicht angelegt“, `.laeuft` des ersten bleibt, der zweite Zweig läuft (`exit=2` des roten Makefiles). Nur dieser Lauf endet.
  (Zwei Fehlversuche meiner Anordnung, ein vererbter Dateideskriptor und bereits vorhandene Refs, lagen nicht am Hook.)

### Stand der Befunde

- **V-1 (Zeitzone): behoben.** `m-p9-ende-lokal` und `m-p9-beginn-lokal` färben genau `p9-zeitzone-ost p9-zeitzone-west`, `p9-beginn-nach-sperre` bleibt frei; der Hook liefert UTC unabhängig von `TZ`.
  Die ganze Gegenprobe am Host ist trotzdem rot (neu V-5).
- **V-2 (root, `p13-ablage-weg`/`p13-mktemp`): offen,** Rückgabe an den Architect unverändert; die Gegenprobe endet als root bewusst mit Exit 2 (Kopfzeile).
- **V-3 (Tabelle in §7 kürzt rote Mengen): nicht neu geprüft** außer den 18 oben, die mit §7 übereinstimmen; Zählung 94 Fälle, 131 Mutationen stimmt.
- **V-4 (Werkzeugliste der README-Zeile): nicht neu geprüft.** Die Kopfzeile der Gegenprobe nennt jetzt `pgrep`, `timeout`.
- **V-5 (neu, mittel): `p9-sigterm-wartet-waise` ist am Host unter Last nicht stabil.** Der Fall wartet nach dem TERM nur `sleep 0.3` und zählt dann per
  `/proc/*/fd/*`-Scan die Halter der Sperre (`haelt`); der Scan über alle Prozesse des Hosts ist teuer (sys-Zeit 70 min, Gesamtlauf 24 statt 5 min).
  In den Einzelläufen (3 bis 18 Mutationen, wenig Parallelität) grün, im Volllauf mit 4 parallelen Gruppen auf einem ausgelasteten Host in 31 bis 42 % der
  Mutationen rot; die unveränderte Kopie blieb grün. Die Gegenprobe läuft laut Spezifikation im Container; wer sie am Host fährt (Review „gezielt selbst“), sieht ROT.
  Behebung (Implementer): auf das Ende der Halter bis zu einer Frist warten statt 0,3 s, oder den Scan auf die Ablage dieser Gruppe begrenzen.
  Ob der Container (5 min 30 s grün, vom Implementer belegt) davon frei ist, habe ich nicht gefahren.

Bewertung: Zeitzonen-Abhängigkeit und die Randformen des Hooks bestätigt, die roten Mengen der 18 Mutanten und der echte Lauf stimmen. Die ganze Gegenprobe am Host
bleibt rot (V-5, Flake), im Container belegt grün. Verifizierbar nach V-5 oder mit Eintrag, dass die Gegenprobe nur im Container zählt.
