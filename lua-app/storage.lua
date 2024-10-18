box.watch('box.status', function()
    if box.info.ro then
        return
    end
    box.schema.create_space('bands', {
        format = {
            { name = 'id', type = 'unsigned' },
            { name = 'bucket_id', type = 'unsigned' },
            { name = 'band_name', type = 'string' },
            { name = 'year', type = 'unsigned' }
        },
        if_not_exists = true
    })
    box.space.bands:create_index('id', { parts = { 'id' }, if_not_exists = true })
    box.space.bands:create_index('bucket_id', { parts = { 'bucket_id' }, unique = false, if_not_exists = true })

    box.schema.create_space('unitInfo', {
        format = {
            { name = 'uid', type = 'string' },
            { name = 'bucket_id', type = 'unsigned' },
            { name = 'band_name', type = 'string' },
            { name = 'year', type = 'unsigned' }
        },
        if_not_exists = true
    })
    box.space.unitInfo:create_index('uid', { parts = { 'uid' }, if_not_exists = true })
    box.space.unitInfo:create_index('bucket_id', {
        parts = { 'bucket_id' }, unique = false, if_not_exists = true
    })

    box.schema.create_space('cache', {
        format = {
            { name = 'hash', type = 'string' },
            { name = 'hash_table', type = 'map' },
            { name = 'bucket_id', type = 'unsigned' },
        },
        if_not_exists = true
    })
    box.space.cache:create_index('hash', { parts = { 'hash' }, if_not_exists = true })
    box.space.cache:create_index('bucket_id', {
        parts = { 'bucket_id' }, unique = false, if_not_exists = true
    })

    box.schema.func.create('get_hashes_by_pattern_storage', { language = 'lua' })
    function get_hashes_by_pattern_storage(pattern)
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

end)