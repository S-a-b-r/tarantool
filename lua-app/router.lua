local vshard = require('vshard')

box.schema.func.create('get_all_hashes_by_pattern', { language = 'lua', if_not_exists = true })
function get_all_hashes_by_pattern(pattern)
    local hashes = {}

    local results = vshard.router.map_callrw('get_hashes_by_pattern', pattern )
    if results ~= nil then
        for _, result in ipairs(results) do
            if result.status == box.OK then
                for _, hash in ipairs(results) do
                    table.insert(hashes, hash)
                end
            end
        end
    end
    return hashes
end


box.schema.func.create('get_hashes_by_pattern', { language = 'lua', if_not_exists = true })
    function get_hashes_by_pattern(pattern)
        local hashes = {}

        local cache = box.space.cache
        -- Loop through all UIDs
        for _, tuple in cache:pairs() do
            -- Check if the UID matches the pattern
            if string.match(tuple[1], pattern) then
                -- Add the UID to the list of matching UIDs
                table.insert(hashes, tuple[1])
            end
        end
        -- Return the list of matching UIDs
        return hashes
    end