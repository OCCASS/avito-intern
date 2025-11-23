CREATE TABLE "user" (
	id        VARCHAR(255) PRIMARY KEY,
	name      VARCHAR(255),
	is_active BOOLEAN DEFAULT false
);

CREATE TABLE pullrequest (
	id         VARCHAR(255) PRIMARY KEY,
	name       VARCHAR(255),
	author_id  VARCHAR(255),
	status     VARCHAR(10) CHECK (status IN ('OPEN', 'MERGED')),
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
	merged_at TIMESTAMP,
	CONSTRAINT fk_author_id
		FOREIGN KEY (author_id)
		REFERENCES "user" (id)
		ON DELETE SET NULL
);

CREATE TABLE pullrequest_reviewer (
	pullrequest_id VARCHAR(255),
	reviewer_id    VARCHAR(255),
	PRIMARY KEY (pullrequest_id, reviewer_id),
	CONSTRAINT fk_pullrequest_id
		FOREIGN KEY (pullrequest_id)
		REFERENCES pullrequest (id)
		ON DELETE CASCADE,
	CONSTRAINT fk_reviewer_id
		FOREIGN KEY (reviewer_id)
		REFERENCES "user" (id)
		ON DELETE CASCADE

);

CREATE TABLE team (
	name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE team_member (
	team_name VARCHAR(255),
	member_id VARCHAR(255),
	PRIMARY KEY (team_name, member_id),
	CONSTRAINT fk_team_name
		FOREIGN KEY (team_name)
		REFERENCES team (name)
		ON DELETE CASCADE,
	CONSTRAINT fk_member_id
		FOREIGN KEY (member_id)
		REFERENCES "user" (id)
		ON DELETE CASCADE
);

