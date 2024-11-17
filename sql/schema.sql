PRAGMA foreign_keys = ON;

-- Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    breadcrumbs TEXT NOT NULL DEFAULT '',
    employment_type INTEGER NOT NULL DEFAULT -1,
    seniority INTEGER NOT NULL DEFAULT -1,
    location TEXT NOT NULL DEFAULT '',
    number_to_hire INTEGER NOT NULL DEFAULT 0,
    experience TEXT NOT NULL DEFAULT '',
    salary TEXT NOT NULL DEFAULT '',
    remote INTEGER NOT NULL DEFAULT -1,
    interview_process TEXT NOT NULL DEFAULT '',
    job_description TEXT NOT NULL DEFAULT '',
    requirements TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_jobs_link ON jobs (link);

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tag TEXT NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_tags_tag ON tags (tag);

-- Create jobs_tags table
CREATE TABLE IF NOT EXISTS jobs_tags (
    job_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (job_id, tag_id)
);

-- Create locations table
CREATE TABLE IF NOT EXISTS locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    address TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    area TEXT NOT NULL DEFAULT '',
    zip_code TEXT NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_locations_address ON locations (address);

-- Create jobs_locations table
CREATE TABLE IF NOT EXISTS jobs_locations (
    job_id INTEGER NOT NULL,
    location_id INTEGER NOT NULL,
    FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
    FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE CASCADE,
    PRIMARY KEY (job_id, location_id)
);
