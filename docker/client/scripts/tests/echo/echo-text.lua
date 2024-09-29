local TEXT = "This is a test of an endpoint that echos its provided string."

request = function()
  return wrk.format("POST", "/api/echoText", nil, TEXT)
end

-- response = function(status, headers, body)
--   if status == 200 then
--     print(body)
--   else
--     print("!!! Unexpected status code: " .. status)
--   end
-- end
