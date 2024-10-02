package com.jtlapp.jvmvsjs.vertxquery;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import io.vertx.core.Future;
import io.vertx.pgclient.PgPool;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;
import io.vertx.sqlclient.Tuple;

import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class SharedQuery {

    private static final ObjectMapper objectMapper = new ObjectMapper();

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
            var paramName = m.group(1);
            if (!paramNames.contains(paramName)) {
                paramNames.add(paramName);
            }
            var paramNumber = paramNames.indexOf(paramName) + 1;
            m.appendReplacement(sb, "\\$" + paramNumber);
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
    public Future<String> executeUsingGson(PgPool db, String jsonArgs) {
        var args = JsonParser.parseString(jsonArgs).getAsJsonObject();
        var tuple = buildTuple(args);

        switch (returnType) {
            case NOTHING:
                // TODO: proper use of preparedQuery?
                return db.preparedQuery(query).execute(tuple).map(res -> "");
            case ROWS:
                return db.preparedQuery(query).execute(tuple).flatMap(this::convertRowsToJson);
            case ROW_COUNT:
                return db.preparedQuery(query).execute(tuple)
                        .map(rowSet -> String.format("{\"rowCount\": %d}", rowSet.rowCount()));
            default:
                throw new SharedQueryException("Unhandled ReturnType");
        }
    }

    private Tuple buildTuple(JsonObject args) {
        Tuple tuple = Tuple.tuple();
        for (var paramName : paramNames) {
            var argPrimitive = args.getAsJsonPrimitive(paramName);
            if (argPrimitive.isNumber()) {
                tuple.addInteger(argPrimitive.getAsInt());
            } else if (argPrimitive.isString()) {
                tuple.addString(argPrimitive.getAsString());
            } else if (argPrimitive.isBoolean()) {
                tuple.addBoolean(argPrimitive.getAsBoolean());
            } else {
                throw new RuntimeException("Unrecognized JSON primitive in "+
                        "parameter '"+ paramName +"'");
            }
        }
        return tuple;
    }

    public Future<String> convertRowsToJson(RowSet<Row> rowSet) {
        var jsonArray = new JsonArray();

        for (Row row : rowSet) {
            var jsonObject = new JsonObject();
            int columnCount = row.size(); // Get the number of columns in the row

            for (int i = 0; i < columnCount; ++i) {
                String columnName = row.getColumnName(i);
                Object value = row.getValue(i); // Get the value, Vert.x abstracts the type

                if (value instanceof String) {
                    jsonObject.addProperty(columnName, (String) value);
                } else if (value instanceof Integer) {
                    jsonObject.addProperty(columnName, (Integer) value);
                } else if (value instanceof Long) {
                    jsonObject.addProperty(columnName, (Long) value);
                } else if (value instanceof Boolean) {
                    jsonObject.addProperty(columnName, (Boolean) value);
                } else if (value instanceof Double) {
                    jsonObject.addProperty(columnName, (Double) value);
                } else if (value instanceof java.time.LocalDate) {
                    jsonObject.addProperty(columnName, value.toString()); // Dates converted to strings
                } else if (value instanceof java.time.LocalDateTime) {
                    jsonObject.addProperty(columnName, value.toString()); // Timestamps converted to strings
                } else if (value == null) {
                    jsonObject.add(columnName, null); // Handle null values
                } else {
                    return Future.failedFuture(
                            new RuntimeException(String.format(
                                    "Column '%s' has unrecognized type '%s'",
                                    columnName, value.getClass().getSimpleName()
                            ))
                    );
                }
            }

            jsonArray.add(jsonObject); // Add each row's JSON object to the array
        }

        return Future.succeededFuture(
                String.format("{\"query\": \"%s\", \"rows\": %s}", name, jsonArray));
    }
}
