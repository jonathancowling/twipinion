package com.infinityworks.cowling.jonathan.twipinion.ingester;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest(
	properties = {
		"config.twitter.bearer="
	}
)
class TweetIngesterApplicationTests {

	@Test
	void contextLoads() {
	}

}
