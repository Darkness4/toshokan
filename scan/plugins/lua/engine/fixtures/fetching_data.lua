local msgpack = require("msgpack")

return msgpack.encode({
  title = "Example Title",
  release_date = 1682390400,   -- Unix timestamp
  tags = { "tag1", "tag2", "tag3" }
})
