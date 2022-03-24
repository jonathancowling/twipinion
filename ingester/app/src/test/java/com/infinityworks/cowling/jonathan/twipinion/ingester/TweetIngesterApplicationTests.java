package com.infinityworks.cowling.jonathan.twipinion.ingester;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest(
	properties = {
		"TWITTER_BEARER="
	}
)
class TweetIngesterApplicationTests {

	@Test
	void contextLoads() {
	}

}
