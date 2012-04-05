#!/usr/bin/python

# utility function, writes inData to inFile
def writeDataToFile(inFile, inData):
	FILE = open(inFile, 'w')
	FILE.writelines(inData)
	FILE.close()

# utility function, returns data from file
def readLinesFromFile(inFile):
	openedFile = open(inFile, "r")
	return openedFile.readlines()
