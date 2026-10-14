**Stand:** gestrichen

Begründung: `slice-v1-abschluss-anmeldung` liefert den Test, den `slice-harness-blackbox-driven` hier der Adresse zuwies: `TestOpenSessionTraegtLesepufferDesAufbaus` wird rot, wenn `Open` eine Session mit neuem Frontend zurückgibt (Mutant U17; Implementer und Verifikation). Die Beobachtung kann in der benannten Form nicht mehr auftreten. Die Grenze (beide Nachrichten in einem Lesevorgang über TCP-Loopback, nicht zugesagt) steht in §6 von `slice-v1-abschluss-anmeldung`. Das Verzeichnis bleibt liegen.
