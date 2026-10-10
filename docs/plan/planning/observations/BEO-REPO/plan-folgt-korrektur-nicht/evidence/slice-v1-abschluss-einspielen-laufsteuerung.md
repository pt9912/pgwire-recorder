**Vorgang:** slice-v1-abschluss-einspielen-laufsteuerung
**Fund:** Der Code-Commit `1b63ef8` änderte den Kommentar des Driving Ports `Player` (Schicht *Ports*) und nannte ihn nur in §3; Kopf (`Berührte Spec-Stellen`) und die Schicht-Abgrenzung in §1 zog erst der Commit des Architect `ee23fad` nach (Review F-567, LOW). `make kopf-check` fing es nicht, weil §1 die Kennung damals nicht nannte.
