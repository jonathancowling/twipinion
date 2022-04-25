package com.infinityworks.cowling.jonathan.twipinion.snowflake;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

import org.springframework.data.repository.CrudRepository;

public interface TweetsRepo extends CrudRepository<TweetsRepo.Tweet, String> {

    @Entity
	@lombok.Data
	@lombok.Builder
	public static class Tweet {
		@Id
		@GeneratedValue(strategy = GenerationType.AUTO)
		String id;
		String data;
	}
}
