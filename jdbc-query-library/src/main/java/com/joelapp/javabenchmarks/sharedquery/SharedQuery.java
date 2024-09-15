package com.joelapp.javabenchmarks.sharedquery;

import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

import java.sql.*;
import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class SharedQuery {

    public enum ReturnType {
        NOTHING("nothing"),
        ROWS("rows"),
        ROW_COUNT("count");

        private final String label;

        ReturnType(String label) {
            this.label = label;
        }

        static ReturnType fromLabel(String label) {
            for (ReturnType candidate : values()) {
                if (candidate.label.equals(label))
                    return candidate;
            }
            throw new RuntimeException(String.format(
                    "ReturnType label '%s' not found", label));
        }
    }

    private final String name;
    private final String query;
    private final ArrayList<String> paramNames = new ArrayList<>();
    private final ReturnType returnType;

    /**
     * Query that is shared across the various benchmarked API servers.
     */
    SharedQuery(String name, String query, ReturnType returnType) {
        this.name = name;

        StringBuilder sb = new StringBuilder();
        Pattern p = Pattern.compile("\\$\\{([a-zA-Z_][a-zA-Z0-9_]*)\\}");
        Matcher m = p.matcher(query);

        while (m.find()) {
            paramNames.add(m.group(1));
            m.appendReplacement(sb, "?");
        }
        m.appendTail(sb);
        this.query = sb.toString();

        this.returnType = returnType;
    }

    /**
     * Returns the name of the shared query.
     */
    public String getName() {
        return name;
    }

    /**
     * Executes the query using arguments in the provided JSON and returns
     * the results in JSON. Parses and generates JSON using GSON.
     * @param db Postgres database to which to connect
     * @param jsonArgs JSON string providing a value for each parameter.
     * @return A JSON representation of the query result. Either an array
     *      of JSON objects or a JSON object having property "rowCount".
     */
    public String executeUsingGson(SharedQueryDB db, String jsonArgs) {

        try (var connection = db.openConnection();
             var statement = connection.prepareStatement(query))
        {
            var args = JsonParser.parseString(jsonArgs).getAsJsonObject();
            supplyArguments(statement, args);

            switch (returnType) {
                case NOTHING:
                    try (var _resultSet = statement.executeQuery()) {
                        return "{}";
                    }
                case ROWS:
                    try (var resultSet = statement.executeQuery()) {
                        return convertResultSetToJson(resultSet);
                    }
                case ROW_COUNT:
                    int count = statement.executeUpdate();
                    return String.format("{\"rowCount\":%d}", count);
                default:
                    throw new SharedQueryException("Unhandled ReturnType");
            }
        }
        catch (SQLException e) {
            throw new SharedQueryException(e.getMessage());
        }
    }

    private void supplyArguments(PreparedStatement statement, JsonObject args)
            throws SQLException
    {
        for (var i = 1; i <= paramNames.size(); ++i) {
            var paramName = paramNames.get(i - 1);
            var argPrimitive = args.getAsJsonPrimitive(paramName);
            if (argPrimitive.isNumber()) {
                statement.setInt(i, argPrimitive.getAsInt());
            } else if (argPrimitive.isString()) {
                statement.setString(i, argPrimitive.getAsString());
            } else if (argPrimitive.isBoolean()) {
                statement.setBoolean(i, argPrimitive.getAsBoolean());
            } else {
                throw new RuntimeException("Unrecognized JSON primitive in "+
                        "parameter '"+ paramName +"'");
            }
        }
    }

    private String convertResultSetToJson(ResultSet resultSet)
            throws SQLException
    {
        var metaData = resultSet.getMetaData();
        var jsonArray = new JsonArray();
        while (resultSet.next()) {
            var jsonObject = new JsonObject();
            for (int i = 1; i <= metaData.getColumnCount(); ++i) {
                var columnName = metaData.getColumnName(i);
                int columnType = metaData.getColumnType(i);
                switch (columnType) {
                    case Types.VARCHAR -> jsonObject.addProperty(
                            columnName, resultSet.getString(i));
                    case Types.INTEGER -> jsonObject.addProperty(
                            columnName, resultSet.getInt(i));
                    case Types.BOOLEAN -> jsonObject.addProperty(
                            columnName, resultSet.getBoolean(i));
                    case Types.DATE -> jsonObject.addProperty(
                            columnName, resultSet.getDate(i).toString());
                    case Types.NUMERIC -> jsonObject.addProperty(
                            columnName, resultSet.getDouble(i));
                    case Types.TIMESTAMP -> jsonObject.addProperty(
                            columnName, resultSet.getTimestamp(i).toString());
                    default -> throw new RuntimeException(String.format(
                            "Column '%s' has unrecognized type name %s",
                            columnName, metaData.getColumnTypeName(i)));
                }
            }
            jsonArray.add(jsonObject);
        }
        return String.format("{\"rows\":%s}", jsonArray);
    }
}
