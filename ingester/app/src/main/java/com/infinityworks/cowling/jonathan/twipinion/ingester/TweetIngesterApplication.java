package com.infinityworks.cowling.jonathan.twipinion.ingester;

import java.time.Clock;
import java.time.Duration;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import com.fasterxml.jackson.core.JsonProcessingException;
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
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;
import org.springframework.util.concurrent.ListenableFuture;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.util.DefaultUriBuilderFactory;
import org.springframework.web.util.UriBuilderFactory;

import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Scheduler;
import reactor.core.scheduler.Schedulers;
import reactor.netty.http.client.HttpClient;

@Slf4j
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
	public Supplier<Mono<TweetResult>> test(
		@Autowired TwitterClient twitter,
		@Autowired ObjectMapper mapper,
		@Autowired KafkaTemplate<String, String> kafka
	) {
		return () -> {
			return twitter.recent()
			    .take(Duration.ofSeconds(15))
				.flatMapIterable((l) -> l)
				.flatMap(tweet -> Mono.fromFuture(() -> {
						ListenableFuture<SendResult<String, String>> l;
						try {
							log.info("test pre send");
							l = kafka.sendDefault(tweet.getId(), mapper.writeValueAsString(tweet));
						} catch (JsonProcessingException e) {
							e.printStackTrace();
							return CompletableFuture.failedFuture(e);
						}
						CompletableFuture<Tweet> f = new CompletableFuture<>() {
							@Override
							public boolean cancel(boolean mayInterruptIfRunning) {
							boolean result = l.cancel(mayInterruptIfRunning);
							super.cancel(mayInterruptIfRunning);
							return result;
							}
						};
						l.addCallback((result) -> {
							log.info("test post send: {}", result);
							f.complete(tweet);
						}, (e) -> {
							log.error("callback: {}", e);
							f.completeExceptionally(e);
						});
						return f;
				})
			).onErrorContinue((e, o) -> {
				log.error("continuing: {}\n{}", o, e);
			}).collect(Collectors.toList())
			.map(tweets -> new TweetResult(tweets));
		};
	}

	@lombok.Value
	public static class TweetResult {
		private List<Tweet> tweets;
	}
}
