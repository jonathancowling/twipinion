package com.infinityworks.cowling.jonathan.twipinion.ingester.client;

import java.net.URI;
import java.time.Clock;
import java.time.Duration;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.util.UriBuilderFactory;

import lombok.extern.slf4j.Slf4j;
import reactor.core.publisher.Flux;
import reactor.core.publisher.SynchronousSink;

@Slf4j
@Service
public class TwitterClient {
    @Autowired
    private WebClient client;

    @Value("${config.twitter.uri}")
    private String uri;

    @Value("${config.twitter.within}")
    private String within;

    @Value("${config.twitter.query}")
    private String query;

    @Value("${config.twitter.bearer}")
    private String token;

    @Value("${config.timeout}")
    private String timeout;

    @Autowired
    private UriBuilderFactory uriBuilder;

    @Autowired
    private Clock clock;

    public Flux<List<Tweet>> recent() {

        String startTime = Instant.now(clock).minus(Duration.parse(within)).toString();

        return Flux.generate(
            () -> uriBuilder.uriString(uri)
                .path("/2/tweets/search/recent")
                .queryParam("tweet.fields", "created_at,entities")
                .queryParam("start_time", startTime)
                .queryParam("query", query)
                .build(),
            (URI recentTweetsUri, SynchronousSink<TweetResponse> sink) -> {
                Optional<TweetResponse> res = client.get()
                    .uri(recentTweetsUri)
                    .header("Authorization", "Bearer " + token)
                    .retrieve()
                    .onStatus(s -> s.isError(), (r) -> r.createException())
                    .bodyToMono(TweetResponse.class)
                    .delayElement(Duration.ofSeconds(1))
                    .blockOptional();

                res.ifPresentOrElse(r -> sink.next(r), () -> sink.complete());

                return res.map(r -> uriBuilder.uriString(uri)
                    .path("/2/tweets/search/recent")
                    .queryParam("tweet.fields", "created_at,entities")
                    .queryParam("start_time", startTime)
                    .queryParam("query", query)
                    .queryParamIfPresent("next_token", Optional.ofNullable(r.getMeta().nextToken))
                    .build())
                    .orElse(null);
        })
            .doOnNext(n -> {
                log.info("got tweet {}", n);
            })
            .doOnError(e -> {
                log.info("failed getting tweet {}", e);
            })
            .flatMap(res -> Flux.just(res.toTweets().toArray((i) -> new Tweet[i])))
            .bufferTimeout(100, Duration.parse(timeout));
    }

    @lombok.Value
    @JsonIgnoreProperties(ignoreUnknown = true)
    private static class TweetResponse {
        Meta meta;
        List<Data> data;

        @lombok.Value
        @JsonIgnoreProperties(ignoreUnknown = true)
        private static class Meta {
            @JsonProperty("result_count")
            int resultCount;
            @JsonProperty("next_token")
            String nextToken;
        }

        @lombok.Value
        @JsonIgnoreProperties(ignoreUnknown = true)
        private static class Data {
            String id;
            String text;
            Entities entities = new Entities();
            @JsonProperty("created_at")
            OffsetDateTime createdAt;

            @lombok.Value
            @JsonIgnoreProperties(ignoreUnknown = true)
            private static class Entities {
                List<Hashtag> hashtags = Collections.emptyList();

                @lombok.Value
                @JsonIgnoreProperties(ignoreUnknown = true)
                private static class Hashtag {
                    String tag;
                }
            }

            private Tweet toTweet() {
                return new Tweet(
                        id,
                        text,
                        entities.getHashtags().stream().map(t -> t.getTag()).collect(Collectors.toList()),
                        createdAt);
            }
        }

        private List<Tweet> toTweets() {
            return data.stream().map(Data::toTweet).collect(Collectors.toList());
        }
    }

    @lombok.Value
    public static class Tweet {
        String id;
        String text;
        List<String> hashtags;
        OffsetDateTime createdAt;

    }

    @lombok.Value
    private class ForbiddenException extends RuntimeException {
        Exception exception;
    }
}
