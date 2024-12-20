box.watch('box.status', function()
    if box.info.ro then
        return
    end

    require('datetime')

    box.schema.create_space('cache', {
        format = {
            { name = 'hash', type = 'string' },
            { name = 'value', type = 'any', is_nullable = true },
            { name = 'bucket_id', type = 'unsigned' },
            { name = 'expired_at', type = 'datetime', is_nullable = true}
        },
        if_not_exists = true
    })
    box.space.cache:create_index('hash', { parts = { 'hash' }, unique = true, if_not_exists = true })
    box.space.cache:create_index('bucket_id', {
        parts = { 'bucket_id' }, unique = false, if_not_exists = true
    })

    box.schema.func.create('get_hashes_by_pattern_storage', { language = 'lua', if_not_exists = true})
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

    box.schema.func.create('drop_expires', { language = 'lua', if_not_exists = true})
    function drop_expires()
        local cache = box.space.cache
        for _, tuple in cache:pairs() do
            if tuple.expired_at > box.time.now() then
                tuple.delete()
            end
        end
        return 'ok'
    end

end)