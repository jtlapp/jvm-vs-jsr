local loggedResponses = {}

function logResponse(statusCode, jsonStr)
    local comboKey = nil

    if jsonStr and jsonStr ~= "" then
        local success, jsonObj = pcall(
            function() return require("lib.dkjson").decode(jsonStr) end
        )

        if success and jsonObj ~= nil then
            local query = jsonObj.query or nil
            local error = jsonObj.error or nil
            comboKey = statusCode .. "|" .. tostring(query) .. "|" .. tostring(error)
        else
          comboKey = statusCode .. "|" .. jsonStr
        end
    else
      comboKey = statusCode .. ""
      jsonStr = "(empty)"
    end

    if not loggedResponses[comboKey] then
      loggedResponses[comboKey] = true
        print(string.format("STATUS %d: %s", statusCode, jsonStr))
    end
end
