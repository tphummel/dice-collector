CREATE TABLE location (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  city TEXT NOT NULL,
  state TEXT NOT NULL
);

CREATE TABLE result (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE shooter (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

CREATE TABLE session (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  locationid INTEGER NOT NULL REFERENCES location (id),
  start TEXT,
  "end" TEXT
);

CREATE TABLE turn (
  sessionid INTEGER NOT NULL REFERENCES session (id),
  shooterid INTEGER NOT NULL REFERENCES shooter (id),
  id INTEGER NOT NULL,
  position TEXT,
  start TEXT,
  "end" TEXT,
  PRIMARY KEY (sessionid, shooterid, id)
);

CREATE TABLE throw (
  sessionid INTEGER NOT NULL,
  shooterid INTEGER NOT NULL,
  turnid INTEGER NOT NULL,
  sequence INTEGER NOT NULL,
  value INTEGER NOT NULL,
  result INTEGER NOT NULL REFERENCES result (id),
  event_prop_comeout TEXT NOT NULL CHECK (event_prop_comeout IN ('Y', 'N')),
  notes TEXT,
  PRIMARY KEY (sessionid, shooterid, turnid, sequence),
  FOREIGN KEY (sessionid, shooterid, turnid) REFERENCES turn (sessionid, shooterid, id)
);
