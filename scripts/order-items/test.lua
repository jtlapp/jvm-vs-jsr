local MAX_USER = 1000
local MAX_ORDER = 4
local SEED = 12345

math.randomseed(SEED)

request = function()
   local user = math.random(1, MAX_USER)
   local order = math.random(1, MAX_ORDER)

   return wrk.format("GET", "/api/select?user=" .. user .. "&order=" .. order)
end
