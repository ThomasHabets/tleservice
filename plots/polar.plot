#!/usr/bin/env gnuplot

set terminal pngcairo
set output "polar.png"
set label 1 at screen 0.02, 0.95 font ":Bold,10"
set label 1 "Satellite XX\nFrom YY"

set label 1 at screen 0.02, 0.95 font ":Bold,10"
set timestamp "Blah blah today's date" offset 1,1
set polar
set angles degrees
set size square
set theta clockwise top
set rrange [90:0]
set style line 12 lc rgb 'grey' lt 1 lw 1
set style line 13 lc rgb 'blue' lt 1 lw 2
set grid polar 30 ls 12
unset ytics
unset xtics
set ttics add ("N" 0, "E" 90, "S" 180, "W" -90) font ":Bold"
set rtics format "%.0f°"
set rtics 30
set rlabel "Altitude" offset -2 font ":Bold"
unset border
set border polar
set style data lines

plot '../pass' w l title '' ls 13
