local SLEEP_MILLIS = 1000

request = function()
  local url = string.format("/api/sleep/%d", SLEEP_MILLIS)
  return wrk.format("GET", url, nil, nil)
end
