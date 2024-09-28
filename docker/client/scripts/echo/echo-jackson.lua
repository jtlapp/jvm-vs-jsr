local JACKSON = "{\"name\":\"John Doe\",\"age\":43.5,\"address\":\"123 Elm St\",\"zip\":62701}"

request = function()
  return wrk.format("POST", "/api/echoJackson", nil, JACKSON)
end

-- response = function(status, headers, body)
--   if status == 200 then
--     print(body)
--   else
--     print("!!! Unexpected status code: " .. status)
--   end
-- end
