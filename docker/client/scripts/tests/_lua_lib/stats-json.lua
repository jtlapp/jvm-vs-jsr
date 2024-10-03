
local function printFloat(name, value, delim)
    print(string.format('  "%s": %f%s', name, value, delim));
end

local function printInt(name, value, delim)
    print(string.format('  "%s": %d%s', name, value, delim));
end

function printStatsJson(summary, latency, requests)
    print("\nJSON {");

    printFloat("latency_min", latency.min, ",");
    printFloat("latency_max", latency.max, ",");
    printFloat("latency_mean", latency.mean, ",");
    printFloat("latency_stdev", latency.stdev, ",");

    print('  "percentiles": [');
    for _, pct in pairs({ 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 99, 99.9, 99.99, 99.999, 100 }) do
        delim = (pct < 100) and "," or ""
        print(string.format("    { \"percentile\": %g, \"microsec_latency\": %g }%s",
            pct, latency:percentile(pct), delim))
    end 
    print("  ],");

    printInt("duration", summary.duration, ",");
    printInt("requests", summary.requests, ",");
    printInt("bytes", summary.bytes, ",");
    printInt("connect_errors", summary.errors.connect, ",");
    printInt("read_errors", summary.errors.read, ",");
    printInt("write_errors", summary.errors.write, ",");
    printInt("status_errors", summary.errors.status, ",");
    printInt("timeout_errors", summary.errors.timeout, "");

    print("}");
end
