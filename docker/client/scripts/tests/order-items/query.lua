package.path = package.path .. ";../?.lua"
require("_lua_lib.response-log")
require("_lua_lib.stats-json")

local MAX_USER = 1000
local MAX_ORDER = 4
local PERCENT_UPDATES = 50
local SEED = 12345

math.randomseed(SEED)

selectRequest = function()
  local userNumber = math.random(1, MAX_USER)
  local orderNumber = math.random(1, MAX_ORDER)

  local orderID = getOrderID(userNumber, orderNumber)
  local postBody = string.format('{"orderID": "%s"}', orderID)
  return wrk.format("POST", "/api/query/orderitems_getOrder", nil, postBody)
end

updateRequest = function()
  local userNumber = math.random(1, MAX_USER)
  local orderNumber = math.random(1, MAX_ORDER)

  local orderID = getOrderID(userNumber, orderNumber)
  local postBody = string.format('{"orderID": "%s"}', orderID)
  return wrk.format("POST", "/api/query/orderitems_boostOrderItems", nil, postBody)
end

request = function()
  if math.random(100) <= PERCENT_UPDATES then
    return updateRequest()
  else
    return selectRequest()
  end
end

response = function(status, headers, body)
  logResponse(status, body)
end

done = function(summary, latency, requests)
  printStatsJson(summary, latency, requests)
end

function getUserID(userNumber)
   return string.format("USER_%06d", userNumber)
end

function getOrderID(userNumber, orderNumber)
   return string.format("%s_ORDER_%d", getUserID(userNumber), orderNumber)
end
