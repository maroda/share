#!/bin/zsh

function cntdots () {
		RN=$1					# first arg = number of runs in CNT seconds
		FAC=8					# defines the "speed" of the iconography
		(( SLP = 1.0 / $FAC ))	# adjust 1-second sleep with factor
		(( ITER = $RN * $FAC ))	# adjust run iterations with factor
		DEC=19967				# decimal unicode: blocks=9631, iching=19967, braille=10495
		CNT=60					# length of each run in approximate seconds

		# displays iconography in descending order of
		# appearance in the Unicode table, also making
		# the cursor invisible for the duration
		#
		tput civis
		for (( z = 1; z <= $ITER; z++ )); do
			for (( c = 0; c < $CNT; c++ )); do
					(( DOT = $DEC - $c ))
					(( OUT = [##16] $DOT ))
					print -n "\\u$OUT"
					sleep $SLP
					print -n "\b"
			done
		done
		print
		tput cvvis
}

print -n "Running for 60 seconds... "
cntdots 1

