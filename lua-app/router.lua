local vshard = require('vshard')

box.schema.func.create('get_all_hashes_by_pattern', { language = 'lua', if_not_exists = true })
function get_all_hashes_by_pattern(pattern)
    local hashes = {}

    local results = vshard.router.map_callrw('get_hashes_by_pattern_storage', pattern )
    for _, result in pairs(results) do
         for _, shard in pairs(result) do
             for _, hash in pairs(shard) do
                table.insert(hashes, hash)
             end
         end
    end
    return hashes
end
