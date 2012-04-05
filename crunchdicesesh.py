#!/usr/bin/python

from PyTom.Dice import data
from PyTom.Dice import logic
from PyTom.Tcrypt.Utilities import io
import sys
import getopt
import string

# main method
if __name__ == "__main__":
	shooters = data.getListOfShooters()
	locations = data.getListOfLocations()
	sessionLocation = ""
	strNextTurnPrompt = "Would you like to enter the next turn? (Y|N): "
	strSession = ""
	currentShooter = ""
	anotherTurn = "Y"
	sessionid = ""
	
	print "Welcome to the Craps Session Entry System\n"
	# choose session location
	for location in locations:
			print str(location[0])+".",location[1]
	strLocationInput = raw_input("Choose Location: ")
	for location in locations:
		if location[0]==int(strLocationInput):
			sessionLocation = location
	# create session record
	sessionid = data.insertSession(sessionLocation[0])
	print "\nSession Created: id=%s" % sessionid
	#turn loop
	while (anotherTurn == "Y"):
		turnid = ""
		print ""
		for shooter in shooters:
			print str(shooter[0])+".",shooter[1]
		
		strShooterInput = raw_input("Choose Shooter: ")
		
		for shooter in shooters:
			if shooter[0] == int(strShooterInput):
				currentShooter = shooter
		
		# create turn record
		turnid = data.insertTurn(sessionid, currentShooter[0])
			
		print "\nTurn Created:", currentShooter[1], "at", sessionLocation[1], "(%s)" % turnid
		strThrows = raw_input("Enter all throws , no spaces, no commas: ")
		# divide throws string into individual chars
		throws = list(strThrows)
		
		# this is where the business happens
		
		finishedThrows = logic.processTurnThrows(throws)
		
		for throw in finishedThrows:
			data.insertThrow(sessionid, currentShooter[0], turnid, throw[0], throw[1], throw[2])
		
		print "Turn Saved. %s: %s" % (currentShooter[1], strThrows)
		
		
		strTurn = currentShooter[1]+": "
		for throw in finishedThrows:
			strTurn += str(throw[0])+","
		strTurn = strTurn.rstrip(",")
		strSession += strTurn+"\n"
		
		# check for next turn
		anotherTurn = raw_input(strNextTurnPrompt).upper()
	
	strSession = strSession.rstrip("\n")
	#io.writeDataToFile(filename,strSession)
	
	print "Session Complete"
	print "Summary--",sessionLocation[1]
	print strSession
	
	
	
