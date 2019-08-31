USE [PremierLeagueTable]
GO

INSERT INTO [dbo].[Seasons]([StartYear], [EndYear])
VALUES(2018, 2019)
GO

INSERT INTO [dbo].[Seasons]([StartYear], [EndYear])
VALUES(2019, 2020)
GO

INSERT INTO [dbo].[Matches] ([MatchCode], [SeasonId], [MatchDay], [HostId], [Host], [GuestId], [Guest], [Date])
VALUES ('MCICHE', 1, 1, 1, 'MCI', 2, 'CHE', '2019-08-10')
GO

INSERT INTO [dbo].[Matches] ([SeasonId], [MatchDay], [HostId], [Host], [GuestId], [Guest], [Date])
VALUES ('LIVEVT', 1, 1, 3, 'LIV', 4, 'EVT', '2019-08-10')
GO

INSERT INTO [dbo].[Matches] ([SeasonId], [MatchDay], [HostId], [Host], [GuestId], [Guest], [Date])
VALUES ('MUTARS', 1, 1, 5, 'MUT', 6, 'ARS', '2019-08-10')
GO

