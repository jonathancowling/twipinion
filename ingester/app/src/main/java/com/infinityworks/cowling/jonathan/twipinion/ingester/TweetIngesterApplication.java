package com.infinityworks.cowling.jonathan.twipinion.ingester;

import java.time.Clock;
import java.util.function.Supplier;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.infinityworks.cowling.jonathan.twipinion.ingester.client.TwitterClient;
import com.infinityworks.cowling.jonathan.twipinion.ingester.client.TwitterClient.Tweet;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.http.client.reactive.ReactorClientHttpConnector;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.util.DefaultUriBuilderFactory;
import org.springframework.web.util.UriBuilderFactory;

import reactor.core.publisher.Flux;
import reactor.netty.http.client.HttpClient;

@SpringBootApplication
public class TweetIngesterApplication {

	public static void main(String[] args) {
		SpringApplication.run(TweetIngesterApplication.class, args);
	}

	@Value("${config.twitter.bearer}")
	String bearer;

	@Value("${config.twitter.query}")
	String query;

	@Bean
	Clock clock() {
		return Clock.systemUTC();
	}

	@Bean
	public ObjectMapper mapper() {
		return new ObjectMapper()
				.disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES)
				.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS)
				.registerModule(new JavaTimeModule());
	}

	@Bean
	public WebClient webClient() {
		return WebClient.builder().clientConnector(
				new ReactorClientHttpConnector(HttpClient.create().proxyWithSystemProperties())).build();
	}

	@Bean
	public UriBuilderFactory uriBuilder() {
		return new DefaultUriBuilderFactory();
	}

	@Bean
	public Supplier<Flux<Tweet>> test(@Autowired TwitterClient twitter, @Autowired ObjectMapper mapper) {
		return () -> {
			return twitter.recent().flatMapIterable((l) -> l);
		};
	}
}
