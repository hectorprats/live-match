using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Threading.Tasks;

namespace SqlBuildingUtilities
{
    class Program
    {
        static async Task Main(string[] args)
        {
            await CreateInsertTeamsSql();
            await CreateInsertSeasonTeamsSql();
        }

        private static async Task CreateInsertTeamsSql()
        {
            const string InsertQuery = @"

INSERT INTO livematches.Teams(Code, Name)
VALUES ({{Code}}, {{Name}});

";

            var path = "Teams/teams.txt";
            var contents = await File.ReadAllLinesAsync(path);
            var teamCodeLine = 21;
            var teamCodes = contents[teamCodeLine].Split('\t');

            var builder = new StringBuilder();

            for (var i = 0; i < 20; i++)
            {
                var insert = InsertQuery.Replace("{{Code}}", $"'{teamCodes[i]}'").Replace("{{Name}}", $"'{contents[i]}'");
                builder.AppendLine(insert);
            }

            var sql = builder.ToString();
        }

        private static async Task CreateInsertSeasonTeamsSql()
        {
            const string InsertQuery = @"

INSERT INTO livematches.SeasonTeams(SeasonId, TeamId)
VALUES (1, {{TeamId}});

";
            var builder = new StringBuilder();

            for (var i = 1; i <= 20; i++)
            {
                var insert = InsertQuery.Replace("{{TeamId}}", $"{i}");
                builder.AppendLine(insert);
            }

            var sql = builder.ToString();
        }
    }
}
