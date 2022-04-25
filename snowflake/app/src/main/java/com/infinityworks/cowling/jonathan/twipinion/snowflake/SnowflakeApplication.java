package com.infinityworks.cowling.jonathan.twipinion.snowflake;

import java.util.function.Consumer;

import org.apache.kafka.streams.kstream.KStream;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.jdbc.core.JdbcTemplate;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@SpringBootApplication
public class SnowflakeApplication {

	public static void main(String[] args) {
		SpringApplication.run(SnowflakeApplication.class, args);
	}

	@Bean
	public Consumer<KStream<String, String>> test(@Autowired JdbcTemplate jdbc, @Autowired TweetsRepo tweets) {
		return (input) -> {
			input.foreach((k, v) -> {
				log.info(k);
				log.info(v);
				tweets.save(TweetsRepo.Tweet.builder().data(v).build());
			});
		};
	}
}
