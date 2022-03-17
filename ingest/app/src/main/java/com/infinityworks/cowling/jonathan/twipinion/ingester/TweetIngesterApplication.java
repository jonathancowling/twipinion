package com.infinityworks.cowling.jonathan.twipinion.ingester;

import java.util.function.Supplier;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

@SpringBootApplication
public class TweetIngesterApplication {

	public static void main(String[] args) {
		SpringApplication.run(TweetIngesterApplication.class, args);
	}

	@Bean
	public Supplier<String> test() {
		return () -> "Hello World";
	}
}
