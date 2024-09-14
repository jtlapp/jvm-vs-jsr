local MAX_USER = 1000
local MAX_ORDER = 4
local SEED = 12345

math.randomseed(SEED)

request = function()
   local userNumber = math.random(1, MAX_USER)
   local orderNumber = math.random(1, MAX_ORDER)

   local orderID = getOrderID(userNumber, orderNumber)
   return wrk.format("GET", "/api/select?order=" ..orderID)
end

function getUserID(userNumber) {
    return string.format("USER_%06d", userNumber)
}

function getOrderID(userNumber, orderNumber) {
    return string.format("%s_ORDER_%06d", getUserID(userNumber), orderNumber)
}
