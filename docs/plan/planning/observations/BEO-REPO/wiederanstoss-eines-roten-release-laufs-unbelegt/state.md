**Stand:** offen

Erstes Auftreten, unter der Schwelle; Ausgang des Risikos „Trockenlauf belegt die Außenwirkung nicht“ (Teil Wiederanstoß) in §6 von `slice-erster-release-freigabe`. Ohne Adresse: Ein Slice dafür wäre ein künstlicher roter Lauf mit Außenwirkung und steht nicht an. Der nächste echte Tag, der rot endet, ist der Beleg oder Gegenbeleg.

Zweites Auftreten (Urteil des Planners, `slice-v1-abschluss-release-v0-1-1`, 2026-10-21; Beleg `evidence/slice-v1-abschluss-release-v0-1-1.md`): Der Lauf des Tags `v0.1.1` endete im Job `release` im ersten Versuch grün; der Wiederanstoß eines roten **Laufs** blieb wieder ungegangen. Der Zähler steht bei 2, unter der Schwelle, weiter ohne Adresse.

**Die andere Hälfte ist seit dem 2026-10-21 belegt:** Der Wiederanlauf eines roten **Jobs** ist dreimal gefahren (Job `homebrew` des Laufs `39551891874`, Versuche 2 bis 4), und zweimal davon ohne Änderung an der Ursache, mit derselben Meldung. Die Beobachtung gilt deshalb nur noch dem roten **Lauf**; die Zusage von `releasing.md` §11 über den Job ist keine Lücke mehr. Ein dritter Beleg für den roten Lauf erreicht die Schwelle.

Kein Auftreten, Vermerk (`slice-erster-release-homebrew-nachweis`, Urteil des Planners, 2026-10-21; kein Zählerzuwachs, der Zähler bleibt bei 2): Der Slice hielt den Wiederanlauf des Jobs `homebrew` und den Nachtrag von Hand als Weg bereit (§1, §6), brauchte beides aber nicht — der Job war im vierten Versuch desselben Laufs grün, und der Tap-Commit ist der Beleg. Der Risiko-Ausgang ist deshalb *entfallen mit Begründung*, und `docs/maintainer/releasing.md` §6 und §11 bekommen den Satz über „Re-run failed jobs“ weiter nicht. Der Eintrag bleibt die Adresse für den nächsten Lauf, der ihn braucht.
