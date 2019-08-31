-- noinspection SqlNoDataSourceInspectionForFile
/*
CREATE DATABASE PremierLeagueTable
GO

USE PremierLeagueTable
GO
*/

CREATE SCHEMA Webhooks

CREATE TABLE Seasons(
    Id SERIAL PRIMARY KEY,
    StartYear INTEGER NOT NULL,
    EndYear INTEGER NOT NULL
)

GO

CREATE TABLE Matches (
    Id SERIAL PRIMARY KEY,
    MatchCode VARCHAR(16),
    SeasonId INT,
    MatchDay INT,
    Host VARCHAR(8),
    Guest VARCHAR(8),
    Date DATE
)
GO

CREATE TABLE MatchResults (
    Id SERIAL PRIMARY KEY,
    MatchId INT,
    MatchCode NVARCHAR(16),
    HostGoals INT,
    GuestGoals INT
)
GO

CREATE TABLE MatchEventType(
    Id SERIAL PRIMARY KEY,
    Name NVARCHAR(32),
    Description NVARCHAR(256)
)
GO

CREATE TABLE MatchEvents (
     Id SERIAL PRIMARY KEY,
     MatchCode VARCHAR(16),
     EventType INT,
     Payload TEXT,
     TimestampCreated TIMESTAMP without time zone default (now() at time zone 'utc')
)
GO

CREATE TABLE Teams (
    Id SERIAL PRIMARY KEY,
    Code NVARCHAR(8),
    Name NVARCHAR(64)
)
GO