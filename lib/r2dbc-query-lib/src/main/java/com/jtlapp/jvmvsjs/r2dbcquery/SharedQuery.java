package com.jtlapp.jvmvsjs.r2dbcquery;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import org.springframework.r2dbc.core.DatabaseClient;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
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
            paramNames.add(paramName);
            m.appendReplacement(sb, ":" + paramName);
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
    public Mono<String> executeUsingGson(DatabaseClient db, String jsonArgs) {
        var args = JsonParser.parseString(jsonArgs).getAsJsonObject();
        var sql = supplyArguments(db.sql(query), args);

        switch (returnType) {
            case NOTHING:
                return sql.then().thenReturn("");
            case ROWS:
                return sql.fetch().all().collectList()
                        .flatMap(this::convertRowsToJson);
            case ROW_COUNT:
                return sql.fetch().rowsUpdated().map(count ->
                        String.format("{\"rowCount\": %d}", count));
            default:
                throw new SharedQueryException("Unhandled ReturnType");
        }
    }

    private DatabaseClient.GenericExecuteSpec supplyArguments(
            DatabaseClient.GenericExecuteSpec sql,
            JsonObject args
    ) {
        for (var paramName : paramNames) {
            var argPrimitive = args.getAsJsonPrimitive(paramName);
            if (argPrimitive.isNumber()) {
                sql = sql.bind(paramName, argPrimitive.getAsInt());
            } else if (argPrimitive.isString()) {
                sql = sql.bind(paramName, argPrimitive.getAsString());
            } else if (argPrimitive.isBoolean()) {
                sql = sql.bind(paramName, argPrimitive.getAsBoolean());
            } else {
                throw new RuntimeException("Unrecognized JSON primitive in "+
                        "parameter '"+ paramName +"'");
            }
        }
        return sql;
    }

    private Mono<String> convertRowsToJson(List<Map<String, Object>> rows) {
        try {
            String jsonArray = objectMapper.writeValueAsString(rows);
            return Mono.just(String.format(
                    "{\"query\": \"%s\", \"rows\": %s}", name, jsonArray));
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }
}
