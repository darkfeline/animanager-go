CREATE TABLE anime (
    aid INTEGER,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    episodecount INTEGER NOT NULL,
    startdate INTEGER,
    enddate INTEGER,
    PRIMARY KEY (aid)
)

CREATE TABLE episode (
    eid INTEGER,
    aid INTEGER NOT NULL,
    type INTEGER NOT NULL,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    length INTEGER NOT NULL,
    user_watched INTEGER NOT NULL CHECK (user_watched IN (0, 1))
        DEFAULT 0,
    PRIMARY KEY (eid),
    FOREIGN KEY (aid) REFERENCES anime (aid)
        ON DELETE CASCADE ON UPDATE CASCADE
)

CREATE TABLE episode_file (
    id INTEGER,
    eid INTEGER NOT NULL,
    path TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (eid) REFERENCES episode (eid)
        ON DELETE CASCADE ON UPDATE CASCADE
)

CREATE TABLE filehash (
    size INTEGER NOT NULL,
    hash TEXT NOT NULL,
    eid INTEGER,
    aid INTEGER,
    filename TEXT,
    UNIQUE(size, hash)
)

CREATE TABLE watching (
    aid INTEGER,
    regexp TEXT NOT NULL,
    offset INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (aid),
    FOREIGN KEY (aid) REFERENCES anime (aid)
	ON DELETE CASCADE ON UPDATE CASCADE
)
