#!/usr/bin/python

def processTurnThrows(throws):
	isComeout = 'Y'
	point = 0
	
	# 0=int value, 1=result, 2=(Y|N) prop_comeout
	
	# converted alpha chars to all numeric
	numThrows = []
	# final throws list. 2d array contain
	outThrows = []
	
	# convert letter values for 10,11,12 to ints
	# handle off and misc rolls, set result also
	
	# INPUT VALUES
	# 2-9 = numbers
	# T = 10
	# E = 11
	# B = 12
	# O = off table
	# M = misc invalid roll
	
	for throw in throws:
		if throw == "T":
			numThrows.append([10,0,0])
		elif throw == "E":
			numThrows.append([11,0,0])
		elif throw == "B":
			numThrows.append([12,0,0])
		elif throw == "O":
			numThrows.append([1,7,0])
		elif throw == "M":
			numThrows.append([1,8,0])
		else:
			numThrows.append([int(throw),0,0])
	
	for throw in numThrows:	
		# set comeout flag for ALL rolls
		throw[2] = isComeout
		# valid throws
		if throw[0] != 1:
			if isComeout == 'Y':
				if throw[0] == 7 or throw[0] == 11:
					throw[1] = 2 #comeout win
				elif throw[0] == 2 or throw[0] == 3 or throw[0] == 12:
					throw[1] = 3 #comeout loss
				else:
					point = throw[0]
					throw[1] = 4 #point
					isComeout = 'N'
			elif isComeout == 'N':
				if throw[0] == point:
					throw[1] = 5 #point win
					isComeout = 'Y'
				elif throw[0] == 7:
					throw[1] = 6 #7out
					isComeout = 'Y'
					point = 0
				else:
					throw[1] = 1 #no consequence
		outThrows.append(throw)
	return outThrows
			

if __name__ == "__main__":
	sampleRolls = [2,3,4,7]
	resultList = processTurnThrows(sampleRolls)
	for throw in resultList:
		print throw[0],throw[1],throw[2]
