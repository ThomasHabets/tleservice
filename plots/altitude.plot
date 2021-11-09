#!/usr/bin/env gnuplot

set terminal png truecolor
set output "altitude.png"

plot 'orbit.dat' using 3 w l title "Altitude in km"
