
	})

	// endpont to redirect shortlinks
	r.GET("/:slug", func(c *gin.Context) {
		slug := c.Param("slug")

		val, err := client.Get(slug).Result()

		if err == redis.Nil {