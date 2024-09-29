package.path = package.path .. ";../?.lua"
require("lib.response-log")

local SEED = 12345
local MAX_ROWS = 1000000
local PERCENT_LONG_REQUESTS = 10
local TAG_CHARS = "0123456789ABCDEF"
local TAG_CHARS_LENGTH = string.len(TAG_CHARS)

math.randomseed(SEED)

long_request = function()
  local tag1 = getRandomTag()
  local tag2 = getRandomTag()

  local postBody = string.format('{"tag1": "%s", "tag2": "%s"}', tag1, tag2)
  return wrk.format("POST", "/api/query/taggedints_sumInts", nil, postBody)
end

short_request = function()
  local id = math.random(1, MAX_ROWS)

  local postBody = string.format('{"id": %d}', id)
  return wrk.format("POST", "/api/query/taggedints_getInt", nil, postBody)
end

request = function()
  if math.random(100) <= PERCENT_LONG_REQUESTS then
    return long_request()
  else
    return short_request()
  end
end

response = function(status, headers, body)
  logResponse(status, body)
end

function getRandomTag()
  local first_char_offset = math.random(1, TAG_CHARS_LENGTH)
  local second_char_offset = math.random(1, TAG_CHARS_LENGTH)
  return TAG_CHARS:sub(first_char_offset, first_char_offset) ..
      TAG_CHARS:sub(second_char_offset, second_char_offset)
end
