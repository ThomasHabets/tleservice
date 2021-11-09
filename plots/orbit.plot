#!/usr/bin/env gnuplot

set terminal png truecolor
set output "orbit.png"

set view equal xyz
set view 80,30,1,1

splot \
      'orbit.dat' using 4:5:6 w l title 'Orbit in space', \
      'orbit.dat' using 7:8:9 w l title 'Orbit relative to ground'
